package servers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/websocket/v2"
	"github.com/magiconair/properties"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/xDefyingGravity/gomcserver"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"watercolormc/internal"
	"watercolormc/internal/app/channels"
	activeServers "watercolormc/internal/app/servers/active"
	"watercolormc/internal/database"
	"watercolormc/internal/utils"
)

type Server struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Port        int    `json:"port"`
	Host        string `json:"host"`
	Version     string `json:"version"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type Versions struct {
	WatercolorVersion string `msgpack:"watercolor"`
	MinecraftVersion  string `msgpack:"minecraft"`
}

type Memory struct {
	Min int `msgpack:"min"`
	Max int `msgpack:"max"`
}

type JavaSettings struct {
	Memory   Memory   `msgpack:"memory"`
	JavaPath string   `msgpack:"javaPath"`
	JvmArgs  []string `msgpack:"jvmArgs"`
}

type ServerConfig struct {
	Versions     Versions     `msgpack:"versions"`
	JavaSettings JavaSettings `msgpack:"javaSettings"`
}

func LoadServerConfig(id string) (*ServerConfig, error) {
	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	if !utils.IsFileExists(serverFolder) {
		return nil, errors.New("server folder not found")
	}

	configFile := filepath.Join(serverFolder, ".watercolor", "config.bin")
	if !utils.IsFileExists(configFile) {
		return nil, errors.New("server config file not found")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := msgpack.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveServerConfig(id string, config *ServerConfig) error {
	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	if !utils.IsFileExists(serverFolder) {
		return errors.New("server folder not found")
	}

	configFile := filepath.Join(serverFolder, ".watercolor", "config.bin")
	data, err := msgpack.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func CreateDefaultServerConfig(minecraftVersion string) *ServerConfig {
	return &ServerConfig{
		Versions: Versions{
			WatercolorVersion: internal.Version,
			MinecraftVersion:  minecraftVersion,
		},
		JavaSettings: JavaSettings{
			Memory: Memory{
				Min: 2048,
				Max: 4096,
			},
			JavaPath: "",
			JvmArgs:  []string{},
		},
	}
}

func makeLogListener(channel, id string) func(string) {
	idCopy := id
	return func(msg string) {
		err := channels.BroadcastToChannel("server:"+channel+":"+idCopy, websocket.TextMessage, []byte(msg))
		zap.L().Debug("broadcasted "+channel, zap.String("id", idCopy), zap.String("message", msg))
		if err != nil {
			zap.L().Error("failed to broadcast "+channel, zap.Error(err))
		}
	}
}

func isRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func secureClone(s string) string {
	return string([]byte(s))
}

func StartServer(id string) error {
	db := database.Get()
	if db == nil {
		return errors.New("database not initialized")
	}

	var (
		name, host, version, createdAt string
		port                           int
	)

	query := `
		SELECT name, port, host, version, created_at
		FROM servers WHERE id = ?`
	err := db.Client.QueryRow(query, id).Scan(&name, &port, &host, &version, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("server not found")
		}
		return err
	}

	stdoutBuf, stderrBuf := bytes.Buffer{}, bytes.Buffer{}
	server := gomcserver.NewServer(name, version)
	server.Directory = utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	server.SetProperty("server-port", strconv.Itoa(port))
	server.SetProperty("server-ip", host)

	s := activeServers.Server{
		ID:           secureClone(id),
		Name:         name,
		Port:         port,
		Host:         host,
		Version:      version,
		CreatedAt:    createdAt,
		StdoutBuffer: &stdoutBuf,
		StderrBuffer: &stderrBuf,
		StdoutWriter: &stdoutBuf,
		StderrWriter: &stderrBuf,
		Instance:     server,
	}
	activeServers.Add(s)

	zap.L().Info("starting server", zap.String("id", id), zap.String("name", name), zap.Int("port", port))
	zap.L().Info("setting stdout listener", zap.String("id", id))

	if err := server.SetEventListener("stdout", makeLogListener("stdout", s.ID)); err != nil {
		zap.L().Error("failed to set stdout listener", zap.Error(err))
		return err
	}
	if err := server.SetEventListener("stderr", makeLogListener("stderr", s.ID)); err != nil {
		zap.L().Error("failed to set stderr listener", zap.Error(err))
		return err
	}

	server.AcceptEULA()

	if err := channels.BroadcastToChannel("server:stdout:"+id, websocket.TextMessage, []byte{}); err != nil {
		return err
	}
	if err := channels.BroadcastToChannel("server:stderr:"+id, websocket.TextMessage, []byte{}); err != nil {
		return err
	}

	if err := channels.BroadcastToChannel("server:stats:"+id, websocket.TextMessage, []byte{}); err != nil {
		return err
	}

	if err := channels.BroadcastToChannel("server:stdin:"+id, websocket.TextMessage, []byte{}); err != nil {
		return err
	}

	if err := channels.BroadcastToChannel("server:players:"+id, websocket.TextMessage, []byte{}); err != nil {
		return err
	}

	channels.SetListener("server:stdin:"+id, func(msg string) error {
		if server.GetPID() == -1 {
			return errors.New("server is not running")
		}
		if err := server.SendCommand(msg); err != nil {
			zap.L().Error("failed to send command to server", zap.Error(err))
			return err
		}
		return nil
	})

	err = server.SetEventListener("playerJoin", func(playerName string, _ int) {
		zap.L().Info("player joined", zap.String("id", s.ID), zap.String("player", playerName))

		err := channels.BroadcastToChannel("server:players:"+s.ID, websocket.TextMessage, []byte("join:"+playerName))
		if err != nil {
			zap.L().Error("failed to broadcast player join", zap.Error(err))
		}

		if internal.SendNotifications {
			err := internal.Notify("A player has joined your server!", playerName+" has joined the server\""+s.Name+"\"!")
			if err != nil {
				panic(err)
			}
		}
	})
	if err != nil {
		return err
	}

	err = server.SetEventListener("playerLeave", func(playerName string, _ int) {
		zap.L().Info("player left", zap.String("id", s.ID), zap.String("player", playerName))

		err := channels.BroadcastToChannel("server:players:"+s.ID, websocket.TextMessage, []byte("leave:"+playerName))
		if err != nil {
			zap.L().Error("failed to broadcast player leave", zap.Error(err))
		}

		if internal.SendNotifications {
			err := internal.Notify("A player has left your server!", playerName+" has left the server \""+s.Name+"\"!")
			if err != nil {
				zap.L().Error("failed to send notification", zap.Error(err))
			}
		}
	})
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(2 * time.Second)
		for {
			time.Sleep(1 * time.Second)

			if !isRunning(server.GetPID()) || server.GetPID() == -1 {
				zap.L().Info("server is offline", zap.String("id", s.ID), zap.Int("pid", server.GetPID()))
				msg, err := json.Marshal(struct {
					Online bool `json:"online"`
				}{Online: false})
				if err != nil {
					zap.L().Error("failed to marshal offline message", zap.Error(err))
					continue
				}

				err = channels.BroadcastToChannel("server:stats:"+s.ID, websocket.TextMessage, msg)
				if err != nil {
					zap.L().Error("failed to broadcast offline stats", zap.Error(err))
					continue
				}

				activeServers.Remove(s.ID)
				channels.RemoveListener("server:stdin:" + id)
				break
			}

			stats, err := server.GetStats()
			if err != nil {
				zap.L().Error("failed to get stats", zap.Error(err))

				if strings.Contains(err.Error(), "no such process") {
					zap.L().Info("server process not found, treating as offline", zap.String("id", s.ID))

					msg, err := json.Marshal(struct {
						Online bool `json:"online"`
					}{Online: false})
					if err != nil {
						zap.L().Error("failed to marshal offline message", zap.Error(err))
						continue
					}

					err = channels.BroadcastToChannel("server:stats:"+s.ID, websocket.TextMessage, msg)
					if err != nil {
						zap.L().Error("failed to broadcast offline stats", zap.Error(err))
						continue
					}

					activeServers.Remove(s.ID)
					channels.RemoveListener("server:stdin:" + id)
					break
				}

				continue
			}

			stats.CPUPercent = stats.CPUPercent / float64(runtime.NumCPU())

			msg, err := json.Marshal(struct {
				Online bool                    `json:"online"`
				Stats  *gomcserver.ServerStats `json:"stats"`
			}{Online: true, Stats: stats})

			if err != nil {
				zap.L().Error("failed to marshal stats message", zap.Error(err))
				continue
			}

			_ = channels.BroadcastToChannel("server:stats:"+s.ID, websocket.TextMessage, msg)
		}
	}()

	config, err := LoadServerConfig(id)
	if err != nil {
		zap.L().Error("failed to load server config", zap.Error(err))
		return err
	}

	err = server.SetMinMemoryMB(config.JavaSettings.Memory.Min)
	if err != nil {
		return err
	}

	err = server.SetMaxMemoryMB(config.JavaSettings.Memory.Max)
	if err != nil {
		return err
	}

	startOpts := &gomcserver.StartOptions{
		StdoutPipe:       &stdoutBuf,
		StderrPipe:       &stderrBuf,
		UseManifestCache: utils.PtrBool(true),
		CacheDir:         utils.PtrString(utils.ExpandHome(internal.WatercolorDirectory + "/cache")),
		JvmOptions:       &config.JavaSettings.JvmArgs,
	}

	if config.JavaSettings.JavaPath != "" {
		startOpts.JavaPath = &config.JavaSettings.JavaPath
	}

	if err := server.Start(startOpts); err != nil {
		activeServers.Remove(id)
		return err
	}

	return nil
}

func StopServer(id string) error {
	server, ok := activeServers.Get(id)
	if !ok {
		return errors.New("server not found")
	}
	if err := server.Instance.SendCommand("stop"); err != nil {
		return err
	}
	activeServers.Remove(id)
	channels.RemoveListener("server:stdin:" + id)
	return nil
}

func GetServerLogs(id string) ([]string, error) {
	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	if !utils.IsFileExists(serverFolder) {
		return nil, errors.New("server logs not found")
	}
	logFile := serverFolder + "/logs/latest.log"
	if !utils.IsFileExists(logFile) {
		return nil, errors.New("server log file not found")
	}
	file, err := os.ReadFile(logFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")
	return lines, nil
}

func getServerFolderAndProperties(id string) (string, *properties.Properties, error) {
	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	if !utils.IsFileExists(serverFolder) {
		err := errors.New("server not found")
		return "", nil, err
	}

	propertiesFile := filepath.Join(serverFolder, "server.properties")
	if !utils.IsFileExists(propertiesFile) {
		err := errors.New("server properties file not found")
		return "", nil, err
	}

	data, err := os.ReadFile(propertiesFile)
	if err != nil {
		return "", nil, err
	}
	props := properties.NewProperties()
	if err = props.Load(data, properties.UTF8); err != nil {
		return "", nil, err
	}
	return serverFolder, props, nil
}
func GetServerWorld(id string) (worldName string, worldPath string, seed string, levelType string, err error) {
	serverFolder, props, err := getServerFolderAndProperties(id)

	worldName = props.GetString("level-name", "world")
	seed = props.GetString("level-seed", "random")
	levelType = props.GetString("level-type", "default")

	worldPath = filepath.Join(serverFolder, worldName)
	if !utils.IsFileExists(worldPath) {
		err = errors.New("world folder not found")
		return
	}

	return
}

func UploadWorld(file multipart.File, serverId string) error {
	serverFolder, props, err := getServerFolderAndProperties(serverId)
	if err != nil {
		return err
	}

	worldName := props.GetString("level-name", "world")
	worldPath := filepath.Join(serverFolder, worldName)

	if utils.IsFileExists(worldPath) {
		if err := os.RemoveAll(worldPath); err != nil {
			return err
		}
	}

	tmpFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	r, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	topLevelFolder := ""
	for _, f := range r.File {
		parts := strings.Split(f.Name, "/")
		if len(parts) > 1 && parts[0] != "__MACOSX" {
			topLevelFolder = parts[0]
			break
		}
	}

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "__MACOSX") {
			continue
		}

		relPath := f.Name
		if topLevelFolder != "" && strings.HasPrefix(f.Name, topLevelFolder+"/") {
			relPath = strings.TrimPrefix(f.Name, topLevelFolder+"/")
		}

		destPath := filepath.Join(worldPath, relPath)

		if !strings.HasPrefix(destPath, filepath.Clean(worldPath)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		_, err = io.Copy(dstFile, rc)
		dstFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetServerProperties(id string) (*properties.Properties, error) {
	_, props, err := getServerFolderAndProperties(id)
	if err != nil {
		return nil, err
	}

	return props, nil
}

func SaveServerProperties(id string, props *properties.Properties) error {
	serverFolder, _, err := getServerFolderAndProperties(id)
	if err != nil {
		return err
	}

	propertiesFile := filepath.Join(serverFolder, "server.properties")
	if !utils.IsFileExists(propertiesFile) {
		return errors.New("server properties file not found")
	}

	data := props.String()

	if err := os.WriteFile(propertiesFile, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}

func InitServer(server Server) error {
	directory := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + server.Id)
	err := utils.CreateIfNotExists(directory)
	if err != nil {
		return err
	}

	watercolorDirectory := filepath.Join(directory, ".watercolor")
	err = utils.CreateIfNotExists(watercolorDirectory)
	if err != nil {
		return err
	}

	config := CreateDefaultServerConfig(server.Version)
	if err := SaveServerConfig(server.Id, config); err != nil {
		return err
	}

	return nil
}

func GetServerPlayers(id string) ([]string, error) {
	server, ok := activeServers.Get(id)
	if !ok {
		return nil, errors.New("server not found")
	}

	return server.Instance.Players, nil
}

func BackupServer(id string) error {
	db := database.Get()
	if db == nil {
		return errors.New("database not initialized")
	}

	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	if !utils.IsFileExists(serverFolder) {
		return errors.New("server folder not found")
	}

	var (
		name, host, version, createdAt string
		port                           int
	)

	query := `
		SELECT name, port, host, version, created_at
		FROM servers WHERE id = ?`
	err := db.Client.QueryRow(query, id).Scan(&name, &port, &host, &version, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("server not found")
		}
		return err
	}

	server := gomcserver.NewServer(name, version)
	server.Directory = serverFolder

	err = server.SendCommand("save-off")
	time.Sleep(2 * time.Second)
	if err != nil {
		return err
	}
	err = server.Backup(true)
	if err != nil {
		return err
	}

	err = server.SendCommand("save-on")
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	return nil
}

func RestoreBackup(serverId string, backup string) error {
	db := database.Get()
	if db == nil {
		return errors.New("database not initialized")
	}

	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)
	if !utils.IsFileExists(serverFolder) {
		return errors.New("server folder not found")
	}

	var (
		name, host, version, createdAt string
		port                           int
	)

	query := `
		SELECT name, port, host, version, created_at
		FROM servers WHERE id = ?`
	err := db.Client.QueryRow(query, serverId).Scan(&name, &port, &host, &version, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("server not found")
		}
		return err
	}

	server := gomcserver.NewServer(name, version)
	server.Directory = serverFolder

	err = server.RestoreBackup(backup)
	if err != nil {
		return err
	}

	return nil
}

func GetServerBackups(serverId string) ([]string, error) {
	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + serverId)
	if !utils.IsFileExists(serverFolder) {
		return nil, errors.New("server folder not found")
	}

	backupDir := filepath.Join(serverFolder, "backups")
	if !utils.IsFileExists(backupDir) {
		return []string{}, nil
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() {
			backups = append(backups, file.Name())
		}
	}

	return backups, nil
}

func DeleteBackup(id string, name string) error {
	serverFolder := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)
	if !utils.IsFileExists(serverFolder) {
		return errors.New("server folder not found")
	}

	backupDir := filepath.Join(serverFolder, "backups")
	if !utils.IsFileExists(backupDir) {
		return errors.New("backup directory not found")
	}

	backupPath := filepath.Join(backupDir, name)
	if !utils.IsFileExists(backupPath) {
		return errors.New("backup file not found")
	}

	if err := os.Remove(backupPath); err != nil {
		return err
	}

	return nil
}
