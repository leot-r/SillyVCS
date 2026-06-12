package models

import (
	"fmt"
	"time"
)

// Committing
type Commit struct {
	Id string `json:"id"`
	File string `json:"file"`
	Time int64 `json:"time"`
	Msg string `json:"msg"`
	Parent string `json:"parent,omitempty"`
	Author string `json:"author"`
}

func NewCommit(id string, file string, msg string, parent string, author string) Commit {
	return Commit {
		Id: id,
		File: file,
		Time: time.Now().Unix(),
		Msg: msg,
		Parent: parent,
		Author: author,
	}
}

type Commits []Commit

// Staging
type StageFile struct {
	FileName string `json:"file_name"`
	IsDir bool `json:"is_dir"`
	Children []StageFile `json:"children"`
}

func (s *StageFile) Walk(fn func(file StageFile)) {
	// USAGE:
	// stagedFiles.Walk(func(file StageFile) {
	//   fmt.Println(file.FileName)
	// })
	fn(*s)
	for _, child := range s.Children {
		child.Walk(fn)
	}
}

func (s *StageFile) Find(name string) *StageFile {
	if s.FileName == name {
		return s
	}
	for i := range s.Children {
		found := s.Children[i].Find(name)
		if found != nil {
			return found
		}
	}
	return nil
}

func (s *StageFile) Print(prefix string, isLast bool) {
	// USAGE:
	// stageFIle.Print("", true)
    connector := "├── "
    if isLast {
        connector = "└── "
    }

    if prefix == "" {
        fmt.Println(s.FileName)
    } else {
        fmt.Println(prefix + connector + s.FileName)
    }

    for i, child := range s.Children {
        isLastChild := i == len(s.Children)-1
        newPrefix := prefix + "│   "
        if isLast {
            newPrefix = prefix + "    "
        }
        child.Print(newPrefix, isLastChild)
    }
}
