package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type StageFile struct {
	FileName string
	IsDir bool
	Children []StageFile
}

func Stage(args []string) {
	path, err := os.Getwd()
	if err != nil {panic(err)}

	var stagedFiles StageFile

	ignored, err := readIgnoreFile(path)
	if err != nil {
		panic(err)
	}

	stagedFiles, err = stageAll(path, ignored)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== STAGING: ===")
	printTree(stagedFiles, "", true)
}

func stageAll(path string, ignored []string) (StageFile, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return StageFile{}, err
	}

	root := StageFile{
		FileName: path,
		IsDir: true,
		Children: []StageFile{},
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
				return StageFile{}, err
			}
			root.Children = append(root.Children, child)
		} else {
			root.Children = append(root.Children, StageFile{
				FileName: file.Name(),
				IsDir: false,
				Children: nil,
			})
		}
	}

	return root, nil
}

func printTree(file StageFile, prefix string, isLast bool) {
    connector := "├── "
    if isLast {
        connector = "└── "
    }

    if prefix == "" {
        fmt.Println(file.FileName)
    } else {
        fmt.Println(prefix + connector + file.FileName)
    }

    for i, child := range file.Children {
        isLastChild := i == len(file.Children)-1
        newPrefix := prefix + "│   "
        if isLast {
            newPrefix = prefix + "    "
        }
        printTree(child, newPrefix, isLastChild)
    }
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
