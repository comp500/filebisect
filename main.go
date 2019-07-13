package main

import (
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
	fmt.Println("Available commands: index, initial, good, bad, help")
}

// Index is the index of all files to bisect
type Index struct {
	TempDirectory string

	Unknown []string
	Ignored []string
	Safe    []string
	Unsafe  []string

	DependencyDefinitions map[string][]string
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

	// Check all indexed files
	var currIndexed []string
	for _, v := range idx.Unknown {
		currIndexed = checkFile(v, currIndexed, idx.TempDirectory)
	}
	for _, v := range idx.Ignored {
		currIndexed = checkFile(v, currIndexed, "")
	}
	for _, v := range idx.Safe {
		currIndexed = checkFile(v, currIndexed, "")
	}
	for _, v := range idx.Unsafe {
		currIndexed = checkFile(v, currIndexed, "")
	}

	// Remove invalid files
	removeInvalidFiles(idx.Unknown, currIndexed)
	removeInvalidFiles(idx.Ignored, currIndexed)
	removeInvalidFiles(idx.Safe, currIndexed)
	removeInvalidFiles(idx.Unsafe, currIndexed)

	// Add new files
	for _, file := range files {
		isIndexed := false
		for _, v := range currIndexed {
			absFile, _ := filepath.Abs(file.Name())
			absValue, _ := filepath.Abs(v)
			if absFile == absValue {
				isIndexed = true
				break
			}
		}
		if !isIndexed {
			idx.Unknown = append(idx.Unknown, file.Name())
		}
	}
}

func checkFile(file string, list []string, tempDir string) []string {
	_, err1 := os.Stat(file)
	if err1 != nil {
		if len(tempDir) > 0 {
			_, err2 := os.Stat(filepath.Join(tempDir, file))
			if err2 != nil {
				fmt.Printf("Error reading file: %v\n", err1)
				return list
			}
		} else {
			fmt.Printf("Error reading file: %v\n", err1)
			return list
		}
	}
	for _, v := range list {
		if v == file {
			fmt.Printf("Error: file %s is duplicated\n", v)
			return list
		}
	}
	return append(list, file)
}

func removeInvalidFiles(list []string, fullList []string) []string {
	newList := list[:0]
	for _, file := range list {
		for _, v := range fullList {
			if file == v {
				newList = append(newList, v)
				break
			}
		}
	}
	return newList
}
