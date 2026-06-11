package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"crypto/sha256"
	"time"

	"SillyVCS/utils"
	"SillyVCS/models"
	"SillyVCS/files"
)


func CommitFile(filePath string) {
	projPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fullPath := filepath.Join(projPath, MetaDirName)

	// Check if repo exists
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
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input: %s\n", data)

	// === Save the meta data ===
	// Make the key for the hash as the file path + the contents of the file
	hashKey := append([]byte(filePath), data...)

	// Get id with SHA-256
	hash := sha256.Sum256(hashKey)

	// Convert it into a hex string
	idHash := fmt.Sprintf("%x", hash)


	// TODO add actually saving the file
	// Write temp file
	files.SaveFile(&data, idHash, filepath.Join(MetaDirName, "snapshots"))
	// Store bytes as a .blob
	// compress??
	// fcync the temp file and rename to .simplevcs/snapshots/.blob


	newCommit := models.Commit{
		Id: idHash,
		File: filePath,
		Time: time.Now().Unix(),
		Msg: "yikers",
		Author: "leo",
	}

	// TODO add temp file, fsync and write later to improve security
	err = files.AddCommit(filepath.Join(MetaDirName, "meta.json"), newCommit)
	if err != nil {
		panic(err)
	}

	// PRINT COMMIT
	commits, err := files.ReadCommits(filepath.Join(MetaDirName, "meta.json"))
	if err != nil {
		fmt.Println("ERROR IN PRINTING COMMITS:")
		fmt.Println(err)
	}
	fmt.Printf("COMMITS:\n  %#v\n", commits)
}
