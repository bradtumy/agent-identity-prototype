package executionlog

import (
	"encoding/json"
	"os"
	"sync"
)

// Logger writes execution results to a file in JSON Lines format.
type Logger struct {
	path string
	mu   sync.Mutex
}

// NewLogger creates a new Logger writing to the specified path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Entry represents a single execution event.
type Entry struct {
	Timestamp string `json:"timestamp"`
	AgentDID  string `json:"agent_did"`
	Role      string `json:"role"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// Log writes the entry to the log file with thread-safety.
func (l *Logger) Log(e Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}
