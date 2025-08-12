package activeServers

import (
	"bytes"
	"io"
	"sync"
)
import (
	"github.com/xDefyingGravity/gomcserver"
)

type Store struct {
	mu      sync.RWMutex
	servers map[string]Server
}

type Server struct {
	ID          string
	Name        string
	Port        int
	Host        string
	Version     string
	Description string
	CreatedAt   string

	StdoutBuffer *bytes.Buffer
	StderrBuffer *bytes.Buffer

	StdoutWriter io.Writer
	StderrWriter io.Writer

	Instance *gomcserver.Server
}

var globalStore = &Store{
	servers: make(map[string]Server),
}

func Get(id string) (Server, bool) {
	globalStore.mu.RLock()
	defer globalStore.mu.RUnlock()

	srv, ok := globalStore.servers[id]
	return srv, ok
}

func Add(srv Server) {
	globalStore.mu.Lock()
	defer globalStore.mu.Unlock()

	globalStore.servers[srv.ID] = srv
}

func Remove(id string) {
	globalStore.mu.Lock()
	defer globalStore.mu.Unlock()

	delete(globalStore.servers, id)
}

func List() []Server {
	globalStore.mu.RLock()
	defer globalStore.mu.RUnlock()

	servers := make([]Server, 0, len(globalStore.servers))
	for _, s := range globalStore.servers {
		servers = append(servers, s)
	}
	return servers
}

func IsOnline(id string) bool {
	globalStore.mu.RLock()
	defer globalStore.mu.RUnlock()

	_, ok := globalStore.servers[id]
	return ok
}
