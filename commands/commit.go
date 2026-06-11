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

// PLAN:
// commit <msg>

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

	// === Save the meta data ===
	// Make the key for the hash as the file path + the contents of the file
	hashKey := append([]byte(filePath), data...)

	// Get id with SHA-256
	hash := sha256.Sum256(hashKey)

	// Convert it into a hex string
	idHash := fmt.Sprintf("%x", hash)

	// TODO compress??
	// fcync the temp file and rename to .simplevcs/snapshots/.blob
	err = files.SaveFile(&data, idHash, filepath.Join(MetaDirName, "snapshots"))
	if err != nil {
		panic(err)
	}

	newCommit := models.NewCommit(idHash, filePath, "yikers", "", "Leo")

	// TODO add temp file, fsync and write later to improve security
	err = files.AddCommit(filepath.Join(MetaDirName, "meta.json"), newCommit)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Committed successfully!\n")
	fmt.Printf("  ID:     %s\n", newCommit.Id)
	fmt.Printf("  File:   %s\n", newCommit.File)
	fmt.Printf("  Msg:    %s\n", newCommit.Msg)
	fmt.Printf("  Author: %s\n", newCommit.Author)
	fmt.Printf("  Parent: %s\n", newCommit.Parent)
	fmt.Printf("  Time:   %s\n", time.Unix(newCommit.Time, 0).Format("02/01/06 - 15:04:05"))
}
