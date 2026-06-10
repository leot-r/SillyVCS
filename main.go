package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Commit struct {
	Id string `json:"id"`
	File string `json:"file"`
	Time int64 `json:"time"`
	Msg string `json:"msg"`
	Parent string `json:"parent,omitempty"`
	Author string `json:"author"`
}

type Commits []Commit

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

func readCommits(path string) (Commits, error){
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Commits{}, nil
		}
		return nil, err
	}

	var commits Commits
	err = json.Unmarshal(data, &commits)
	return commits, err
}

func writeCommits(path string, commits Commits) error {
	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func addCommit(path string, commit Commit) error {
	commits, err := readCommits(path)
	if err != nil {
		return err
	}
	
	commits = append(commits, commit)

	return writeCommits(path, commits)
}

func checkIfInit(rootPath string) (exists bool, error error) {
	// Check if a repo has been initalized in the current directory
	_, err := os.Stat(filepath.Join(rootPath, metaDirName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
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
	
	exists, err := checkIfInit(rootPath)
	if err != nil {
		panic(err)
	}
	if exists == true {
		fmt.Println("Repo is already initialized in this dir")
		return nil
	}

	// If file doesnt exist setup
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
	// Read bytes
	data, err := os.ReadFile(filePath)
	fmt.Printf("Input: %s\n", data)

	// Make the key for the hash as the file path + the contents of the file
	hashKey := append([]byte(filePath), data...)

	// Get id with SHA-256
	hash := sha256.Sum256(hashKey)

	newCommit := Commit{
		Id: fmt.Sprintf("%x", hash),
		File: filePath,
		Time: time.Now().Unix(),
		Msg: "yikers",
		Author: "leo",
	}

	// TODO add temp file, fsync and write later to improve security
	err = addCommit(filepath.Join(metaDirName, "meta.json"), newCommit)
	if err != nil {
		panic(err)
	}

	// PRINT COMMIT
	commits, err := readCommits(filepath.Join(metaDirName, "meta.json"))
	if err != nil {
		fmt.Println("ERROR IN PRINTING COMMITS:")
		fmt.Println(err)
	}
	fmt.Printf("COMMITS:\n  %#v\n", commits)
}
