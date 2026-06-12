package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"SillyVCS/utils"
	"SillyVCS/models"
	"SillyVCS/files"
)

func Stage(args []string) {
	projPath, err := os.Getwd()
	if err != nil {panic(err)}

	// Path to staging area / single file
	stagingPath := filepath.Join(projPath, args[2])

	// Check for initialized repo
	repoExists, err := utils.CheckIfInit(filepath.Join(projPath, MetaDirName))
	if err != nil {
		panic(err)
	}
	if repoExists == false {
		fmt.Println("There is not an initialized repo int this directory :/")
		return
	}

	var stagedFiles models.StageFile

	// Find what files/dirs are in the .vcsignore
	ignored, err := readIgnoreFile(projPath)
	if err != nil {
		panic(err)
	}

	// Find what files are being staged
	stagedFiles, err = stageAll(stagingPath, ignored)
	if err != nil {
		panic(err)
	}

	fmt.Println("Staging...")
	stagedFiles.Print("", true)

	// Add stage files to stage.json
	err = files.WriteStageFile(filepath.Join(projPath, MetaDirName, "staging.json"), stagedFiles)
	if err != nil {
		panic(err)
	}
}

func stageAll(path string, ignored []string) (models.StageFile, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return models.StageFile{}, err
	}

	root := models.StageFile{
		FileName: path,
		IsDir: true,
		Children: []models.StageFile{},
	}

	// Recursive bs
	for _, file := range dir {
		if isIgnored(file.Name(), ignored) {
			continue
		}
		if file.IsDir() {
			// recurse ino subdirectory here
			child, err := stageAll(filepath.Join(path, file.Name()), ignored)
			if err != nil {
				return models.StageFile{}, err
			}
			root.Children = append(root.Children, child)
		} else {
			root.Children = append(root.Children, models.StageFile{
				FileName: filepath.Join(path, file.Name()),
				IsDir: false,
				Children: nil,
			})
		}
	}

	return root, nil
}

func readIgnoreFile(path string) ([]string, error) {
	// Check For .vscignore
	data, err := os.ReadFile(filepath.Join(path, ".vcsignore"))
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	ignored := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			ignored = append(ignored, line)
		}
	}

	ignored = append(ignored, ".sillyvcs")
	
	return ignored, nil
}

func isIgnored(name string, ignored []string) bool {
	for _, pattern := range ignored {
		if name == pattern {
			return true
		}
	}
	return false
}
