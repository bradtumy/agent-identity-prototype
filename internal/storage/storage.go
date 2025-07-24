package storage

import (
	"encoding/json"
	"os"
	"sync"
)

// Agent stores agent identity information on disk.
type Agent struct {
	DID        string                 `json:"did"`
	Owner      string                 `json:"owner"`
	Metadata   map[string]interface{} `json:"metadata"`
	Credential interface{}            `json:"credential"`
}

// FileStore stores agents to a JSON file.
type FileStore struct {
	path string
	mu   sync.Mutex
	data map[string]Agent
}

// NewFileStore creates a file backed store at path.
func NewFileStore(path string) *FileStore {
	fs := &FileStore{path: path, data: map[string]Agent{}}
	fs.load()
	return fs
}

func (fs *FileStore) load() {
	b, err := os.ReadFile(fs.path)
	if err != nil {
		return
	}
	json.Unmarshal(b, &fs.data)
}

func (fs *FileStore) save() error {
	b, err := json.MarshalIndent(fs.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.path, b, 0644)
}

// Save stores an agent record.
func (fs *FileStore) Save(a Agent) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.data[a.DID] = a
	return fs.save()
}
