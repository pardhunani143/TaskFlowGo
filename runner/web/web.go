package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pardhunani143/TaskFlowGo/runner/types"
)

type Server struct {
	listenAddr string
	processor  types.TaskProcessor
}

func StartHttpServer(listenAddr string, proc types.TaskProcessor) {

	server := Server{
		listenAddr: listenAddr,
		processor:  proc,
	}

	server.setupRoutes()
	log.Printf("Starting HTTP server on %s", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

}

func (s *Server) setupRoutes() {
	http.HandleFunc("/task", s.HandleTask)
	http.HandleFunc("/health", s.handleHealth)
}

func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received task request")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200! - OK"))
}

func (s *Server) HandleTask(w http.ResponseWriter, r *http.Request) {
	var task types.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate task
	if task.Type == "" || task.ID == "" {
		http.Error(w, "invalid task: missing type or id", http.StatusBadRequest)
		return
	}

	// Submit task to processor
	s.processor.Submit(task)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"task_id": task.ID,
		"status":  "accepted",
	})
}
