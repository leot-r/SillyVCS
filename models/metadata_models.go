package models

import (
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
