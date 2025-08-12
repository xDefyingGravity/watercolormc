package internal

import (
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"path/filepath"
)

type Settings struct {
	BasePath string `msgpack:"base_path"`
}

func (s *Settings) GetBasePath() string {
	if s.BasePath == "" {
		return WatercolorDirectory
	}
	return s.BasePath
}

func (s *Settings) SetBasePath(basePath string) {
	s.BasePath = basePath
}

func (s *Settings) Save() error {
	data, err := msgpack.Marshal(s)
	if err != nil {
		return err
	}

	filePath := WatercolorDefaultDirectory + "/settings.bin"
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func LoadSettings() (*Settings, error) {
	filePath := WatercolorDefaultDirectory + "/settings.bin"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &Settings{BasePath: WatercolorDirectory}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var settings Settings
	if err := msgpack.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
