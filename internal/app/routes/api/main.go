package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/magiconair/properties"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"watercolormc/internal"
	"watercolormc/internal/app/servers"
	activeServers "watercolormc/internal/app/servers/active"
	"watercolormc/internal/database"
	"watercolormc/internal/paper/plugins"
	"watercolormc/internal/utils"
)

func RegisterApiRoutes(app *fiber.App) {
	app.Get("/api/upstatus", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/api/servers", func(c *fiber.Ctx) error {
		db := database.Get()

		rows, err := db.Client.Query(`
		SELECT id, name, port, host, version, description, created_at
		FROM servers
	`)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("error querying servers")
		}
		defer rows.Close()

		var servers []map[string]interface{}

		for rows.Next() {
			var (
				id          string
				name        string
				port        int
				host        string
				version     string
				description string
				createdAt   time.Time
			)

			if err := rows.Scan(&id, &name, &port, &host, &version, &description, &createdAt); err != nil {
				zap.L().Error("error scanning server row", zap.Error(err))
				return c.Status(fiber.StatusInternalServerError).SendString("error scanning server row")
			}

			status := "offline"
			if activeServers.IsOnline(id) {
				status = "online"
			}

			servers = append(servers, map[string]interface{}{
				"id":          id,
				"name":        name,
				"port":        port,
				"host":        host,
				"version":     version,
				"description": description,
				"createdAt":   createdAt.Format(time.RFC3339),
				"status":      status,
			})
		}

		if err := rows.Err(); err != nil {
			zap.L().Error("error iterating over server rows", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error iterating over server rows")
		}

		if servers == nil {
			servers = []map[string]interface{}{}
		}
		zap.L().Info("retrieved servers", zap.Int("count", len(servers)))
		return c.JSON(servers)
	})

	app.Post("/api/servers", func(c *fiber.Ctx) error {
		db := database.Get()

		var server servers.Server

		if err := c.BodyParser(&server); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			rawBody := c.Body()
			zap.L().Error("error parsing request body",
				zap.Error(err),
				zap.ByteString("raw_body", rawBody),
				zap.String("content_type", c.Get("Content-Type")),
			)
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		processed, processError := utils.PreprocessVersion(server.Version)
		if processError != nil {
			zap.L().Error("error preprocessing version", zap.Error(processError))
			return c.Status(fiber.StatusBadRequest).SendString("invalid version format")
		}
		server.Version = processed

		if server.Id == "" || server.Name == "" || server.Port <= 0 || server.Host == "" || server.Version == "" {
			zap.L().Error("missing required fields in server data", zap.Any("server", server))
			return c.Status(fiber.StatusBadRequest).SendString("missing required fields")
		}

		_, err := db.Client.Exec(`
		INSERT INTO servers (id, name, port, host, version, description)
		VALUES (?, ?, ?, ?, ?, ?)
	`, server.Id, server.Name, server.Port, server.Host, server.Version, server.Description)

		if err != nil {
			zap.L().Error("error inserting server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error inserting server")
		}

		err = servers.InitServer(server)
		if err != nil {
			zap.L().Error("error initializing server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error initializing server")
		}

		return c.JSON(map[string]interface{}{
			"id":          server.Id,
			"name":        server.Name,
			"port":        server.Port,
			"host":        server.Host,
			"version":     server.Version,
			"description": server.Description,
			"createdAt":   time.Now().Format(time.RFC3339),
		})
	})

	app.Delete("/api/servers/:id", func(c *fiber.Ctx) error {
		db := database.Get()

		id := c.Params("id")

		_, err := db.Client.Exec(`
		DELETE FROM servers WHERE id = ?
	`, id)

		serverPath := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id)

		if utils.IsFileExists(serverPath) {
			if err := os.RemoveAll(serverPath); err != nil {
				zap.L().Error("error removing server directory", zap.Error(err))
				return c.Status(fiber.StatusInternalServerError).SendString("error removing server directory")
			}
		}

		if err != nil {
			zap.L().Error("error deleting server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error deleting server")
		}

		return c.SendString("ok")
	})

	app.Get("/api/ipinfo", func(c *fiber.Ctx) error {
		privateIp, err := utils.GetPrivateIP()
		if err != nil {
			zap.L().Error("error getting private IP", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting private IP")
		}
		publicIp, err := utils.GetPublicIP()
		if err != nil {
			zap.L().Error("error getting public IP", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting public IP")
		}

		return c.JSON(map[string]string{
			"privateIp": privateIp,
			"publicIp":  publicIp,
		})
	})

	app.Post("/api/servers/start/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		err := servers.StartServer(id)
		if err != nil {
			zap.L().Error("error starting server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error starting server")
		}

		return c.SendString("ok")
	})

	app.Post("/api/servers/stop/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		err := servers.StopServer(id)
		if err != nil {
			zap.L().Error("error stopping server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error stopping server")
		}

		return c.SendString("ok")
	})

	app.Get("/api/servers/logs/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		logs, err := servers.GetServerLogs(id)
		if err != nil {
			zap.L().Error("error getting server logs", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting server logs")
		}

		return c.JSON(logs)
	})

	app.Get("/api/servers/:id/world", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		worldName, worldPath, worldSeed, worldType, err := servers.GetServerWorld(id)
		if err != nil {
			zap.L().Error("error getting server world path", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting server world path")
		}

		return c.JSON(map[string]string{"name": worldName, "path": worldPath, "seed": worldSeed, "type": worldType})
	})

	app.Post("/api/servers/:id/world/upload", func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "could not open uploaded file")
		}
		defer file.Close()

		if err := servers.UploadWorld(file, c.Params("id")); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to upload world: "+err.Error())
		}

		return c.SendString("ok")
	})

	app.Get("/api/servers/:id/properties", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		props, err := servers.GetServerProperties(id)
		if err != nil {
			zap.L().Error("error getting server properties", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting server properties")
		}

		return c.JSON(props.Map())
	})
	app.Post("/api/servers/:id/properties", func(c *fiber.Ctx) error {
		db := database.Get()
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		var newProps map[string]string
		if err := c.BodyParser(&newProps); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid JSON body")
		}

		propsFile := utils.ExpandHome(internal.WatercolorDirectory + "/servers/" + id + "/server.properties")
		props := properties.NewProperties()

		for key, value := range newProps {
			if key == "server-port" {
				port, err := strconv.Atoi(value)
				if err != nil || port <= 0 || port > 65535 {
					return c.Status(fiber.StatusBadRequest).SendString("invalid server-port value")
				}
				_, err = db.Client.Exec("UPDATE servers SET port = ? WHERE id = ?", port, id)
				if err != nil {
					zap.L().Error("failed to update server port in database", zap.Error(err))
					return c.Status(fiber.StatusInternalServerError).SendString("failed to update server port")
				}
			} else if key == "server-ip" {
				if value != "" && !utils.IsValidIP(value) {
					return c.Status(fiber.StatusBadRequest).SendString("invalid server-ip value")
				}

				_, err := db.Client.Exec("UPDATE servers SET host = ? WHERE id = ?", value, id)
				if err != nil {
					zap.L().Error("failed to update server host in database", zap.Error(err))
					return c.Status(fiber.StatusInternalServerError).SendString("failed to update server host")
				}
			}

			err := props.SetValue(key, value)
			if err != nil {
				return err
			}
		}

		f, err := os.Create(propsFile)
		if err != nil {
			zap.L().Error("failed to open properties file for writing", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("failed to write properties")
		}
		defer f.Close()

		if _, err := props.WriteComment(f, "# updated via api", properties.UTF8); err != nil {
			zap.L().Error("failed to write properties", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("failed to write properties")
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/api/servers/:id/config", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		config, err := servers.LoadServerConfig(id)
		if err != nil {
			zap.L().Error("error loading server config", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error loading server config")
		}
		return c.JSON(config)
	})

	app.Post("/api/servers/:id/config", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		var config servers.ServerConfig
		if err := c.BodyParser(&config); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if err := servers.SaveServerConfig(id, &config); err != nil {
			zap.L().Error("error saving server config", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error saving server config")
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/api/servers/:id/players", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		players, err := servers.GetServerPlayers(id)
		if err != nil {
			zap.L().Error("error getting server players", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting server players")
		}

		return c.JSON(players)
	})

	app.Get("/api/minecraft/uuid/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")

		resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + username)
		if err != nil {
			zap.L().Error("failed to fetch mojang uuid", zap.Error(err))
			return c.Status(500).SendString("failed to fetch mojang uuid")
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			zap.L().Error("failed to fetch mojang uuid", zap.Int("status", resp.StatusCode))
			return c.Status(resp.StatusCode).SendString("mojang user not found")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.L().Error("failed to read mojang response body", zap.Error(err))
			return c.Status(500).SendString("failed to read response")
		}

		return c.Type("application/json").Send(body)
	})

	app.Post("/api/servers/:id/backup", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		err := servers.BackupServer(id)
		if err != nil {
			zap.L().Error("error creating server backup", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error creating server backup")
		}

		return c.SendString("ok")
	})

	app.Post("/api/servers/:id/restore", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		var request struct {
			BackupName string `json:"backupName"`
		}

		if err := c.BodyParser(&request); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if request.BackupName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("backup name is required")
		}

		err := servers.RestoreBackup(id, request.BackupName)
		if err != nil {
			zap.L().Error("error restoring server backup", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error restoring server backup")
		}

		return c.SendString("ok")
	})

	app.Get("/api/servers/:id/backups", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		backups, err := servers.GetServerBackups(id)
		if err != nil {
			zap.L().Error("error getting server backups", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting server backups")
		}

		return c.JSON(backups)
	})

	app.Delete("/api/servers/:id/backups", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		var request struct {
			BackupName string `json:"backupName"`
		}

		if err := c.BodyParser(&request); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if request.BackupName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("backup name is required")
		}

		err := servers.DeleteBackup(id, request.BackupName)
		if err != nil {
			zap.L().Error("error deleting server backup", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error deleting server backup")
		}

		return c.SendString("ok")
	})

	app.Post("/api/servers/:id/plugins", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		var request struct {
			Plugins []string `json:"plugins"`
		}

		if err := c.BodyParser(&request); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if len(request.Plugins) == 0 {
			return c.Status(fiber.StatusBadRequest).SendString("no plugins provided")
		}

		err := plugins.AddMultipleToServer(id, request.Plugins)
		if err != nil {
			zap.L().Error("error adding plugins to server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error adding plugins to server")
		}

		return c.SendString("ok")
	})

	app.Get("/api/servers/:id/plugins", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		installedPlugins, err := plugins.ListPlugins(id)
		if err != nil {
			zap.L().Error("error listing plugins for server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error listing plugins] for server")
		}

		return c.JSON(installedPlugins)
	})

	app.Delete("/api/servers/:id/plugins/:pluginName", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		pluginName := c.Params("pluginName")
		if pluginName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing plugin name")
		}

		err := plugins.RemoveFromServer(id, pluginName)
		if err != nil {
			zap.L().Error("error removing plugin from server", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error removing plugin from server")
		}

		return c.SendString("ok")
	})

	app.Get("/api/servers/:id/plugins/manifest", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		manifest, err := plugins.GetServerPluginsFromManifest(id)
		if err != nil {
			zap.L().Error("error getting plugin manifest", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error getting plugin manifest")
		}
		if manifest == nil {
			return c.Status(fiber.StatusNotFound).SendString("no plugin manifest found for this server")
		}
		return c.JSON(manifest)
	})

	app.Post("/api/servers/:id/plugins/manifest", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		var newPlugin plugins.Plugin

		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, c.Body(), "", "  ")
		if err != nil {
			fmt.Println("invalid json:", err)
		} else {
			fmt.Println(prettyJSON.String())
		}

		if err := c.BodyParser(&newPlugin); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if err := plugins.AddPluginToManifest(id, newPlugin.Id, newPlugin.JarName); err != nil {
			zap.L().Error("error adding plugin to manifest", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error adding plugin to manifest")
		}

		return c.SendString("ok")
	})

	app.Delete("/api/servers/:id/plugins/manifest/:pluginId", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing server ID")
		}

		pluginId := c.Params("pluginId")
		if pluginId == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing plugin ID")
		}

		if err := plugins.RemovePluginFromManifest(id, pluginId); err != nil {
			zap.L().Error("error removing plugin from manifest", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error removing plugin from manifest")
		}

		return c.SendString("ok")
	})

	app.Get("/api/settings", func(c *fiber.Ctx) error {
		settings, err := internal.LoadSettings()
		if err != nil {
			zap.L().Error("error loading settings", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error loading settings")
		}

		if settings == nil {
			settings = &internal.Settings{
				BasePath: internal.WatercolorDirectory,
			}
			if err := settings.Save(); err != nil {
				zap.L().Error("failed to save default settings", zap.Error(err))
				return c.Status(fiber.StatusInternalServerError).SendString("failed to save default settings")
			}
		}

		return c.JSON(settings)
	})

	app.Post("/api/settings", func(c *fiber.Ctx) error {
		var newSettings internal.Settings
		if err := c.BodyParser(&newSettings); err != nil {
			zap.L().Error("error parsing request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if newSettings.BasePath == "" {
			newSettings.BasePath = internal.WatercolorDirectory
		}

		if err := newSettings.Save(); err != nil {
			zap.L().Error("error saving settings", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("error saving settings")
		}

		internal.WatercolorDirectory = newSettings.BasePath

		return c.SendString("ok")
	})
}
