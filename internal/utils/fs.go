package utils

import "os"

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreateIfNotExists(folderPath string) error {
	if IsFileExists(folderPath) {
		return nil
	}

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func RemoveIfExists(path string) error {
	if !IsFileExists(path) {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func ExpandHome(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return homeDir + path[1:]
	}
	return path
}