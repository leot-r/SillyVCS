package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"SillyVCS/files"
	"SillyVCS/models"
	"SillyVCS/utils"
)


var metaDirName string = ".sillyvcs"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("run <help> for commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		fmt.Print("Commands:\n  init <>\n  commit <>\n")
	case "init":
		if len(os.Args) < 3 {
			Init("")
		} else {
			Init(os.Args[2])
		}
	case "commit":
		if len(os.Args) < 3 {
			fmt.Println("Please supply file to commit")
			return
		}
		CommitFile(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

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

	fullPath := filepath.Join(rootDir, metaDirName)
	
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
	err = SetupMetadata(rootPath)
	if err != nil {
		panic(err)
	}
	return nil
}

func SetupMetadata(rootPath string) error {
	// Create .sillyvcs DIR
	err := os.Mkdir(filepath.Join(rootPath, metaDirName), 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println(".sillyvcs dir already exists")
		} else {
			return err
		}
	}
	fmt.Println("Created .sillyvcs dir")

	// Create .sillyvcs/snapshots DIR
	err = os.Mkdir(filepath.Join(rootPath, metaDirName, "snapshots"), 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println(".sillyvcs/snapshots dir already exists")
		} else {
			return err
		}
	}
	fmt.Println("Created .sillyvcs/snapshots dir")

	// Create metadata file
	file, err := os.Create(filepath.Join(rootPath, metaDirName, "meta.json"))
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("meta file already exists")
		} else {
			return err
		}
	} else {
		// Initialize meta file with empty array
		_, err := file.WriteString("[]")
		if err != nil {
			return err
		}
	}
	fmt.Println("meta file created")
	defer file.Close()

	return nil
}

func CommitFile(filePath string) {
	projPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fullPath := filepath.Join(projPath, metaDirName)

	// TODO add init checking
	repoExists, err := utils.CheckIfInit(fullPath)
	if err != nil {
		panic(err)
	}
	if repoExists == false {
		fmt.Println("There is not an initialized repo int this directory :/")
		return
	}


	// Read bytes
	data, err := os.ReadFile(filePath)
	fmt.Printf("Input: %s\n", data)

	// Make the key for the hash as the file path + the contents of the file
	hashKey := append([]byte(filePath), data...)

	// Get id with SHA-256
	hash := sha256.Sum256(hashKey)

	newCommit := models.Commit{
		Id: fmt.Sprintf("%x", hash),
		File: filePath,
		Time: time.Now().Unix(),
		Msg: "yikers",
		Author: "leo",
	}

	// TODO add temp file, fsync and write later to improve security
	err = files.AddCommit(filepath.Join(metaDirName, "meta.json"), newCommit)
	if err != nil {
		panic(err)
	}

	// PRINT COMMIT
	commits, err := files.ReadCommits(filepath.Join(metaDirName, "meta.json"))
	if err != nil {
		fmt.Println("ERROR IN PRINTING COMMITS:")
		fmt.Println(err)
	}
	fmt.Printf("COMMITS:\n  %#v\n", commits)
}
