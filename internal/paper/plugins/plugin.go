package plugins

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"watercolormc/internal"
	"watercolormc/internal/utils"
)

type Plugin struct {
	Id      string `msgpack:"id" json:"id"`
	JarName string `msgpack:"jar_name" json:"jar_name"`
}

type PluginManifest struct {
	Plugins []Plugin `msgpack:"plugins" json:"plugins"`
}

func AddPluginToManifest(serverId string, pluginId string, pluginJarPath string) error {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return errors.New("server not found")
	}

	manifestPath := filepath.Join(serverPath, "plugins.bin")

	var manifest PluginManifest
	if utils.IsFileExists(manifestPath) {
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("failed to read manifest: %w", err)
		}
		err = msgpack.Unmarshal(data, &manifest)
		if err != nil {
			return fmt.Errorf("failed to unmarshal manifest: %w", err)
		}
	}

	for _, p := range manifest.Plugins {
		if p.Id == pluginId {
			return fmt.Errorf("plugin %s already exists in server %s", pluginId, serverId)
		}
	}

	manifest.Plugins = append(manifest.Plugins, Plugin{Id: pluginId, JarName: pluginJarPath})
	data, err := msgpack.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

func RemovePluginFromManifest(serverId string, pluginId string) error {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return errors.New("server not found")
	}

	manifestPath := filepath.Join(serverPath, "plugins.bin")

	if !utils.IsFileExists(manifestPath) {
		return fmt.Errorf("manifest not found for server %s", serverId)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	err = msgpack.Unmarshal(data, &manifest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	for i, p := range manifest.Plugins {
		if p.Id == pluginId || p.JarName == pluginId {
			manifest.Plugins = append(manifest.Plugins[:i], manifest.Plugins[i+1:]...)
		}
	}
	data, err = msgpack.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	return nil
}

// GetServerPluginsFromManifest GetServerPlugins retrieves the list of plugin IDs for a given server ID
func GetServerPluginsFromManifest(serverId string) ([]Plugin, error) {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return nil, errors.New("server not found")
	}

	manifestPath := filepath.Join(serverPath, "plugins.bin")

	if !utils.IsFileExists(manifestPath) {
		return nil, fmt.Errorf("manifest not found for server %s", serverId)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	err = msgpack.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	return manifest.Plugins, nil
}

func GetServerPluginFromJarName(serverId string, jarName string) (*Plugin, error) {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return nil, errors.New("server not found")
	}

	manifestPath := filepath.Join(serverPath, "plugins.bin")

	if !utils.IsFileExists(manifestPath) {
		return nil, fmt.Errorf("manifest not found for server %s", serverId)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	err = msgpack.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	for _, plugin := range manifest.Plugins {
		if plugin.JarName == jarName {
			return &plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin %s not found in server %s", jarName, serverId)
}

func AddToServer(serverId string, pluginUrl string) error {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return errors.New("server not found")
	}

	pluginsPath := filepath.Join(serverPath, "plugins")

	err := utils.CreateIfNotExists(pluginsPath)
	if err != nil {
		return err
	}

	resp, err := http.Get(pluginUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download plugin: status %s", resp.Status)
	}

	parsedUrl, err := url.Parse(pluginUrl)
	if err != nil {
		return err
	}

	filename := path.Base(parsedUrl.Path)
	if filename == "" || filename == "." || filename == "/" {
		return errors.New("could not determine filename from URL")
	}

	outputPath := filepath.Join(pluginsPath, filename)
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func AddMultipleToServer(serverId string, pluginUrls []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(pluginUrls))

	for _, pluginUrl := range pluginUrls {
		wg.Add(1)
		go func(pluginUrl string) {
			defer wg.Done()
			if err := AddToServer(serverId, pluginUrl); err != nil {
				errCh <- fmt.Errorf("failed to add %s: %w", pluginUrl, err)
			}
		}(pluginUrl)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) == 0 {
		return nil
	}

	var allErrs []string
	for err := range errCh {
		allErrs = append(allErrs, err.Error())
	}

	return errors.New(strings.Join(allErrs, "\n"))
}

func RemoveFromServer(serverId string, pluginName string) error {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return errors.New("server not found")
	}

	pluginsPath := filepath.Join(serverPath, "plugins")

	pluginPath := filepath.Join(pluginsPath, pluginName)
	if !utils.IsFileExists(pluginPath) {
		return fmt.Errorf("plugin %s not found in server %s", pluginName, serverId)
	}

	err := os.Remove(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to remove plugin %s: %w", pluginName, err)
	}

	return nil
}

func ListPlugins(serverId string) ([]string, error) {
	serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)

	if !utils.IsFileExists(serverPath) {
		return nil, errors.New("server not found")
	}

	pluginsPath := filepath.Join(serverPath, "plugins")

	if !utils.IsFileExists(pluginsPath) {
		return nil, fmt.Errorf("no plugins found for server %s", serverId)
	}

	files, err := os.ReadDir(pluginsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %w", err)
	}

	var plugins []string
	for _, file := range files {
		if !file.IsDir() {
			plugins = append(plugins, file.Name())
		}
	}

	return plugins, nil
}
