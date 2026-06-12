package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"SillyVCS/utils"
)


const MetaDirName string = ".sillyvcs"


func Init(rootDir string) error {
	// If a path is not supplied use where the command was executed
	rootPath := rootDir
	if rootDir == "" {
		var err error
		rootPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	fullPath := filepath.Join(rootDir, MetaDirName)
	
	exists, err := utils.CheckIfInit(fullPath)
	if err != nil {
		panic(err)
	}
	if exists == true {
		fmt.Println("Repo is already initialized in this dir")
		return nil
	}

	// If file doesnt exist setup
	fmt.Printf("Initializing a repo in %s...\n", rootDir)
	err = setupMetadata(rootPath)
	if err != nil {
		panic(err)
	}
	return nil
}


func setupMetadata(rootPath string) error {
	// Create .sillyvcs DIR
	err := os.Mkdir(filepath.Join(rootPath, MetaDirName), 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println(".sillyvcs dir already exists")
		} else {
			return err
		}
	}
	fmt.Println("Created .sillyvcs dir")

	// Create .sillyvcs/snapshots DIR
	err = os.Mkdir(filepath.Join(rootPath, MetaDirName, "snapshots"), 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println(".sillyvcs/snapshots dir already exists")
		} else {
			return err
		}
	}
	fmt.Println("Created .sillyvcs/snapshots dir")

	// Create metadata file
	file, err := os.Create(filepath.Join(rootPath, MetaDirName, "meta.json"))
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("meta file already exists")
		} else {
			return err
		}
	} else {
		// Initialize meta file with empty array
		_, err := file.WriteString("")
		if err != nil {
			return err
		}
	}
	fmt.Println("meta file created")
	defer file.Close()

	// Create staging metadata file
	stagingFile, err := os.Create(filepath.Join(rootPath, MetaDirName, "staging.json"))
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("staging file already exists")
		} else {
			return err
		}
	} else {
		_, err := file.WriteString("[]")
		if err != nil {
			return err
		}
	}
	fmt.Println("staging file created")
	defer stagingFile.Close()

	return nil
}
