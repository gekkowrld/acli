package config

import (
	"log"
	"os"
)

// Create a directory and overwrite anything
// if overwrite is set
func CreateDir(dirName string, overwrite bool) bool {
	isExists, err := PathExists(dirName)
	if err != nil {
		log.Fatal(err)
	}

	if !overwrite && FileExists(dirName) && isExists {
		log.Fatalf("\"%s\" exists and is a file\n", dirName)
	} else if DirExists(dirName) && !overwrite && isExists {
		log.Fatalf("\"%s\" exists and is a directory\n", dirName)
	} else if overwrite && FileExists(dirName) && isExists {
		err = os.Remove(dirName)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func CreateFile(filename string, content string, overwrite bool) bool {
	// Check if the path exists first
	isExist, err := PathExists(filename)
	if err != nil {
		log.Fatal("Error: %v\n", err)
	}

	// If the path isn't there now create the file
	// The check is to ensure that no accidental loss of data as this
	// is a write operation not append.

	if !overwrite && isExist && FileExists(filename) {
		log.Fatalf("\"%s\" already exists and is a file\n", filename)
	} else if !overwrite && isExist && DirExists(filename) {
		log.Fatalf("\"%s\" already exists and is a directory\n", filename)
	}

	// If a directory and overwrite is passed, then remove all the directories

	if overwrite && isExist && DirExists(filename) {
		err = os.RemoveAll(filename)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return true
}
