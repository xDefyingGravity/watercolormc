package utils

import (
	"strings"
	"watercolormc/internal/paper"
)

func PreprocessVersion(version string) (string, error) {
	if strings.HasPrefix(version, "paper-") {
		version = strings.TrimPrefix(version, "paper-")

		url, err := paper.GetDownloadURL(version)
		if err != nil {
			return "", err
		}
		return url, nil
	}

	return version, nil
}
