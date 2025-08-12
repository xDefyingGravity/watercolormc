package paper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetDownloadURL(version string) (string, error) {
	type VersionInfo struct {
		Builds []int `json:"builds"`
	}

	resp, err := http.Get("https://api.papermc.io/v2/projects/paper/versions/" + version)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var info VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}

	if len(info.Builds) == 0 {
		return "", fmt.Errorf("no builds found for version %s", version)
	}

	latestBuild := info.Builds[len(info.Builds)-1]
	url := fmt.Sprintf(
		"https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar",
		version, latestBuild, version, latestBuild,
	)

	return url, nil
}
