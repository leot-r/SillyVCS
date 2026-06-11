package files

import (
	"os"
	"fmt"
	"path/filepath"
)


func SaveFile(data *[]byte, idHash string, snapshotDirPath string) error {
	tmpFilePath := filepath.Join(snapshotDirPath, fmt.Sprintf("%s.tmp", idHash))
	finalFilePath := filepath.Join(snapshotDirPath, fmt.Sprintf("%s.blob", idHash))

	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}

	_, err = tmpFile.Write(*data)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFilePath)
		return err
	}

	tmpFile.Sync() // fsync
	tmpFile.Close()

	return os.Rename(tmpFilePath, finalFilePath)
}
