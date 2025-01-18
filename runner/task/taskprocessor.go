package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pardhunani143/TaskFlowGo/runner/types"
)

type Processor struct {
	tasks    chan types.Task
	statuses map[string]*taskStatus // Map
	handlers map[types.TaskType]TaskHandler
	mu       sync.RWMutex
	config   types.RunnerConfig
}

type taskStatus struct {
	Status    types.TaskStatus
	Error     string
	UpdatedAt time.Time
}

type TaskHandler interface {
	Execute(task *types.Task) error
}

func NewProcessor(workers int, config types.RunnerConfig) *Processor {
	p := &Processor{
		tasks:    make(chan types.Task, 100),
		statuses: make(map[string]*taskStatus),
		handlers: make(map[types.TaskType]TaskHandler),
		config:   config,
	}

	// Register handlers
	p.handlers[types.TaskTypeGo] = &GoHandler{}
	p.handlers[types.TaskTypePrometheus] = &PrometheusHandler{}
	p.handlers[types.TaskTypeShell] = &ShellHandler{}

	// Start worker pool
	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

func (p *Processor) Submit(task types.Task) error {
	p.tasks <- task
	p.updateStatus(task.ID, types.StatusPending, "")
	return nil
}

func (p *Processor) GetStatus(taskID string) (types.TaskStatus, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if status, exists := p.statuses[taskID]; exists {
		return status.Status, nil
	}
	return "", fmt.Errorf("task not found")
}

func (p *Processor) worker() {
	for task := range p.tasks {
		handler, exists := p.handlers[task.Type]
		if !exists {
			p.updateStatus(task.ID, types.StatusFailed, "unknown task type")
			continue
		}

		p.updateStatus(task.ID, types.StatusRunning, "")
		err := handler.Execute(&task)
		if err != nil {
			log.Printf("Task %s failed: %v", task.ID, err)
			p.updateStatus(task.ID, types.StatusFailed, err.Error())
		} else {
			p.updateStatus(task.ID, types.StatusCompleted, "")
		}
	}
}

func (p *Processor) updateStatus(taskID string, status types.TaskStatus, errMsg string) {
	p.mu.Lock()
	p.statuses[taskID] = &taskStatus{
		Status:    status,
		Error:     errMsg,
		UpdatedAt: time.Now(),
	}
	p.mu.Unlock()

	log.Println("Task", taskID, "status:", status)

	// Report to manager asynchronously
	go p.reportToManager(taskID, status, errMsg)
}

func (p *Processor) reportToManager(taskID string, status types.TaskStatus, errMsg string) {

	log.Println("Reporting task status to manager")
	log.Println("p.config.ManagerURL", p.config.ManagerURL)
	payload := types.TaskStatusUpdate{
		TaskID:    taskID,
		Status:    status,
		Error:     errMsg,
		Timestamp: time.Now(),
	}

	_, err := http.Post(p.config.ManagerURL+"/status", "application/json",
		bytes.NewBuffer(mustJSON(payload)))
	if err != nil {
		log.Printf("Failed to report status to manager: %v", err)
	}
}

func mustJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
