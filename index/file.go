package index

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// File is a file that has been indexed
type File struct {
	Status         string   `toml:"status"`
	StatusOriginal string   `toml:"status-original,omitempty"`
	BadCount       int      `toml:"bad-count,omitzero"`
	Dependencies   []string `toml:"dependencies"`

	fileName     string
	currLocation string
	isInTemp     bool
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
			file.isInTemp = true
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
		file.isInTemp = false
	}
	return true
}

// MoveToTemp moves the file to the temporary directory
func (file *File) MoveToTemp(tempDir string) error {
	newPath := filepath.Join(tempDir, file.fileName)
	if !file.isInTemp {
		err := os.Rename(file.currLocation, newPath)
		if err != nil {
			err = moveFile(file.currLocation, newPath)
			if err != nil {
				return err
			}
		}
		file.currLocation = newPath
		file.isInTemp = true
	}
	return nil
}

// MoveToTest moves the file to the testing directory
func (file *File) MoveToTest() error {
	if file.isInTemp {
		err := os.Rename(file.currLocation, file.fileName)
		if err != nil {
			err = moveFile(file.currLocation, file.fileName)
			if err != nil {
				return err
			}
		}
		file.currLocation = file.fileName
		file.isInTemp = false
	}
	return nil
}

func moveFile(source, destination string) (err error) {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	fi, err := src.Stat()
	if err != nil {
		return err
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm
	dst, err := os.OpenFile(destination, flag, perm)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		dst.Close()
		os.Remove(destination)
		return err
	}
	err = dst.Close()
	if err != nil {
		return err
	}
	err = src.Close()
	if err != nil {
		return err
	}
	err = os.Remove(source)
	if err != nil {
		return err
	}
	return nil
}
