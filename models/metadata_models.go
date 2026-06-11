package models

type Commit struct {
	Id string `json:"id"`
	File string `json:"file"`
	Time int64 `json:"time"`
	Msg string `json:"msg"`
	Parent string `json:"parent,omitempty"`
	Author string `json:"author"`
}

type Commits []Commit
