package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// TODO: change format to embed dependencies in with files, should look better

const indexName = "file-bisect-index.toml"

func main() {
	var idx Index
	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		idx = Index{}
	} else {
		idx, err = LoadIndex()
		if err != nil {
			fmt.Printf("Error loading index file: %v\n", err)
			os.Exit(1)
		}
	}

	if len(os.Args) != 2 {
		fmt.Println("1 argument must be specified!")
		printHelp()
		os.Exit(1)
	}

	switch strings.ToLower(os.Args[1]) {
	case "index":
		idx.Init()
		idx.Refresh()
		err := idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	case "help":
		printHelp()
	default:
		fmt.Println("Invalid command!")
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Available commands: index, good, bad, help")
}

// File is a file that has been indexed
type File struct {
	Status         string   `toml:"status"`
	StatusOriginal string   `toml:"status-original,omitempty"`
	BadCount       int      `toml:"bad-count,omitzero"`
	Dependencies   []string `toml:"dependencies"`

	fileName     string
	currLocation string
}

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

// Init checks that the file is correct, and stores the file name
func (file *File) Init(fileName string) error {
	file.fileName = fileName
	switch file.Status {
	case "":
		file.Status = "unknown"
	// Valid statuses
	case "unknown":
	case "good":
	case "bad":
	case "ignored":
	case "errored":
	// Not one of the above:
	default:
		return errors.New("invalid file status for " + fileName)
	}
	return nil
}

// Check checks that the file exists and is valid
func (file *File) Check(tempDir string) bool {
	_, err1 := os.Stat(file.fileName)
	if err1 != nil {
		if len(tempDir) > 0 && (file.Status == "" || file.Status == "unknown" || file.Status == "bad") {
			_, err2 := os.Stat(filepath.Join(tempDir, file.fileName))
			if err2 != nil {
				fmt.Printf("Error reading file: %v\n", err1)
				return false
			}
			currLoc, err := filepath.Abs(filepath.Join(tempDir, file.fileName))
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return false
			}
			file.currLocation = currLoc
		} else {
			fmt.Printf("Error reading file: %v\n", err1)
			return false
		}
	} else {
		currLoc, err := filepath.Abs(file.fileName)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return false
		}
		file.currLocation = currLoc
	}
	return true
}
