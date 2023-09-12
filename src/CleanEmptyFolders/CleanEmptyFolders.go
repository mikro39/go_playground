package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == os.EOF {
		return true, nil
	}
	return false, err
}

func cleanupEmptyFolders(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// if it's a directory and it's empty, remove it
		if info.IsDir() {
			empty, err := isDirEmpty(path)
			if err != nil {
				return err
			}
			if empty {
				fmt.Println("Removing:", path)
				return os.Remove(path)
			}
		}
		return nil
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cleanup <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	if err := cleanupEmptyFolders(dir); err != nil {
		fmt.Printf("Error cleaning up empty folders: %s\n", err)
		os.Exit(1)
	}
}

// To run this, type "$ go run cleanup.go /path/to/directory"