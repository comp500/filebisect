package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// Index is the index of all files to bisect
type Index struct {
	TempDirectory string `toml:"temp-directory"`

	Files map[string]File
}

// LoadIndex loads the index file
func LoadIndex() (Index, error) {
	var idx Index
	if _, err := toml.DecodeFile(indexName, &idx); err != nil {
		return Index{}, err
	}
	return idx, nil
}

// Save saves the index file
func (idx Index) Save() error {
	f, err := os.Create(indexName)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	// Disable indentation
	enc.Indent = ""
	return enc.Encode(idx)
}

// Init checks all the current files, to see if they are valid
func (idx *Index) Init() {
	if idx.Files == nil {
		idx.Files = make(map[string]File)
	}

	// Init and check all files
	for k, v := range idx.Files {
		err := v.Init(k)
		if err != nil {
			fmt.Printf("Error reading index: %v\n", err)
		}
		if !v.Check(idx.TempDirectory) {
			v.StatusOriginal = v.Status
			v.Status = "errored"
		}
		idx.Files[k] = v
	}
}

// Refresh reads the current folder and updates all the files/folders in it
func (idx *Index) Refresh() {
	// Create a temp directory if it doesn't exist
	if _, err := os.Stat(idx.TempDirectory); os.IsNotExist(err) {
		dir, err := ioutil.TempDir("", "file-bisect-")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
		}
		idx.TempDirectory = dir
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	// Add new files
	for _, file := range files {
		// Ignore file-bisect-index.toml
		if strings.Contains(file.Name(), indexName) {
			continue
		}
		// Ignore folders, they break all the assumptions of this code
		if file.IsDir() {
			continue
		}
		absFile, err := filepath.Abs(file.Name())
		if err != nil {
			continue
		}

		isIndexed := false
		for _, v := range idx.Files {
			if absFile == v.currLocation {
				isIndexed = true
				break
			}
		}
		if !isIndexed {
			newFile := File{}
			newFile.Init(file.Name())
			idx.Files[file.Name()] = newFile
		}
	}
}

// Good states that the files in the main directory are good
func (idx *Index) Good() {
	// Mark all non-temp files as good
	for k, v := range idx.Files {
		if v.Status == "errored" || v.Status == "ignored" {
			continue
		}
		if !v.isInTemp {
			v.Status = "good"
			v.BadCount = 0
		}
		idx.Files[k] = v
	}

	idx.moveFiles()
}

// Bad states that the files in the main directory are bad, one of them is very bad
func (idx *Index) Bad() {
	// Increase badcount
	numInDir := 0
	kTest := ""
	for k, v := range idx.Files {
		if v.Status != "unknown" {
			continue
		}
		if !v.isInTemp {
			v.BadCount++
			numInDir++
			kTest = k
		}
		idx.Files[k] = v
	}
	// If there's only one, it MUST be the offender
	if numInDir == 1 {
		badFile := idx.Files[kTest]
		badFile.Status = "bad"
		idx.Files[kTest] = badFile
	}

	idx.moveFiles()
}

func (idx *Index) moveFiles() {
	// Get unknown files, sort by goodness, then by randomness
	var unknownFiles []*File
	for _, v := range idx.Files {
		// The pointer to v doesn't change, we need to alloc a new var to get a new pointer for each
		newFile := v
		if v.Status == "unknown" {
			unknownFiles = append(unknownFiles, &newFile)
		}
	}
	if len(unknownFiles) == 0 {
		fmt.Println("Done!")
		return
	}
	sortByGoodness(unknownFiles)

	// Move half of the files to the temp folder
	var testHalf []*File
	var tempHalf []*File
	if len(unknownFiles) == 1 {
		testHalf = unknownFiles
		tempHalf = unknownFiles[:0]
	} else {
		testHalf = unknownFiles[:len(unknownFiles)/2]
		tempHalf = unknownFiles[len(unknownFiles)/2:]
	}
	for _, v := range testHalf {
		err := v.MoveToTest()
		if err != nil {
			fmt.Printf("Error moving file: %v\n", err)
		}
	}
	for _, v := range tempHalf {
		err := v.MoveToTemp(idx.TempDirectory)
		if err != nil {
			fmt.Printf("Error moving file: %v\n", err)
		}
	}

	// Put data back into main map
	for _, v := range unknownFiles {
		idx.Files[v.fileName] = *v
	}
}

// Sort the list, first by goodness, then by randomness
func sortByGoodness(files []*File) {
	// Sort randomly
	for i := range files {
		j := rand.Intn(i + 1)
		files[i], files[j] = files[j], files[i]
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].BadCount < files[j].BadCount
	})
}

// Ignore ignores a file, it's that simple!
func (idx *Index) Ignore(file string) {
	fileStr, ok := idx.Files[file]
	if !ok {
		fmt.Println("Couldn't find that file!")
		return
	}
	fileStr.Status = "ignored"
	idx.Files[file] = fileStr
}
