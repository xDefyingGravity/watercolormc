package channels

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	channelsMu       sync.RWMutex
	channels         = make(map[string]map[*websocket.Conn]bool)
	channelListeners = make(map[string]func(string) error)
)

// Init initializes the websocket route and handles connections
func Init(app *fiber.App) {
	app.Get("/channels/:channel", websocket.New(func(c *websocket.Conn) {
		channel := c.Params("channel")
		if channel == "" {
			_ = c.Close()
			return
		}

		channelsMu.Lock()
		if channels[channel] == nil {
			channels[channel] = make(map[*websocket.Conn]bool)
		}
		channels[channel][c] = true
		channelsMu.Unlock()

		defer func() {
			channelsMu.Lock()
			delete(channels[channel], c)
			if len(channels[channel]) == 0 {
				delete(channels, channel)
			}
			channelsMu.Unlock()
			_ = c.Close()
		}()

		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			channelsMu.RLock()
			listener, hasListener := channelListeners[channel]
			channelsMu.RUnlock()

			if hasListener {
				go func(m []byte) {
					if err := listener(string(m)); err != nil {
						fmt.Printf("channel listener error: %v\n", err)
					}
				}(msg)
				continue
			}

			channelsMu.RLock()
			for conn := range channels[channel] {
				if conn != c {
					_ = conn.WriteMessage(msgType, msg)
				}
			}
			channelsMu.RUnlock()
		}
	}))
}

// BroadcastToChannel sends a message to all connections in a channel
func BroadcastToChannel(channel string, msgType int, msg []byte) error {
	channelsMu.Lock()
	if channels[channel] == nil {
		channels[channel] = make(map[*websocket.Conn]bool)
	}
	conns := channels[channel]
	channelsMu.Unlock()

	channelsMu.RLock()
	defer channelsMu.RUnlock()

	for conn := range conns {
		if err := conn.WriteMessage(msgType, msg); err != nil {
			return fmt.Errorf("failed to send message to channel %s: %w", channel, err)
		}
	}

	return nil
}

// SetListener assigns a listener function for a channel
func SetListener(channel string, listener func(string) error) {
	channelsMu.Lock()
	defer channelsMu.Unlock()
	channelListeners[channel] = listener
}

// RemoveListener removes the listener for a channel
func RemoveListener(channel string) {
	channelsMu.Lock()
	defer channelsMu.Unlock()
	delete(channelListeners, channel)
}

// Cleanup closes and removes dead connections and cleans empty channels
func Cleanup() {
	channelsMu.Lock()
	defer channelsMu.Unlock()

	for channel, conns := range channels {
		for conn := range conns {
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				_ = conn.Close()
				delete(conns, conn)
			}
		}
		if len(conns) == 0 {
			delete(channels, channel)
		}
	}
}