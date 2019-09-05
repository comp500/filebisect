package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const indexName = "file-bisect-index.toml"

func main() {
	rand.Seed(time.Now().Unix())

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

	if len(os.Args) < 2 {
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
	case "ignore":
		if len(os.Args) != 3 {
			fmt.Println("The file to ignore must be specified!")
			printHelp()
			os.Exit(1)
		}
		idx.Init()
		idx.Ignore(filepath.Clean(os.Args[2]))
		err := idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	case "good":
		idx.Init()
		idx.Good()
		err := idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	case "bad":
		idx.Init()
		idx.Bad()
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
	fmt.Println("Usage: filebisect [command]")
	fmt.Println("Available commands: index, ignore, good, bad, help")
}
