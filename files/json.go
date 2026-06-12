package files

import (
	"encoding/json"
	"os"

	"SillyVCS/models"
)

// Commiting
func ReadCommits(path string) (models.Commits, error){
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return models.Commits{}, nil
		}
		return nil, err
	}

	var commits models.Commits
	err = json.Unmarshal(data, &commits)
	return commits, err
}

func WriteCommits(path string, commits models.Commits) error {
	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddCommit(path string, commit models.Commit) error {
	commits, err := ReadCommits(path)
	if err != nil {
		return err
	}
	
	commits = append(commits, commit)
	return WriteCommits(path, commits)
}

// Staging
func ReadStageFile(StageFilePath string) (models.StageFile, error) {
	data, err := os.ReadFile(StageFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.StageFile{}, nil
		}
		return models.StageFile{}, err
	}

	var stageFile models.StageFile
	err = json.Unmarshal(data, &stageFile)
	return stageFile, err
}

func WriteStageFile(path string, stageFile models.StageFile) error {
	data, err := json.MarshalIndent(stageFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ClearStageFile(path string) error {
	return os.WriteFile(path, []byte("{}"), 0644)
}
