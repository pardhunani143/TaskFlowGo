package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/pardhunani143/TaskFlowGo/runner/types"
)

type Manager struct {
	runners map[string]types.RunnerConfig
	tasks   map[string]types.Task
	mu      sync.RWMutex
}

func main() {
	m := &Manager{
		runners: make(map[string]types.RunnerConfig),
		tasks:   make(map[string]types.Task),
	}

	// Setup routes
	http.HandleFunc("/register", m.handleRegister)
	http.HandleFunc("/status", m.handleStatus)
	http.HandleFunc("/heartbeat", m.handleHeartbeat)

	// Start server
	log.Println("Starting manager on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func (m *Manager) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config types.RunnerConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	m.runners[config.RunnerID] = config
	m.mu.Unlock()

	log.Printf("Runner registered: %s", config.RunnerID)
	w.WriteHeader(http.StatusOK)
}

func (m *Manager) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var update types.TaskStatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Task %s status: %s", update.TaskID, update.Status)
	w.WriteHeader(http.StatusOK)
}

func (m *Manager) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var heartbeat struct {
		RunnerID string `json:"runner_id"`
		Status   string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&heartbeat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Heartbeat from %s: %s", heartbeat.RunnerID, heartbeat.Status)
	w.WriteHeader(http.StatusOK)
}
