// Config represents the runner configuration
package types

import "time"

const (
	ActionUpdateConfig TaskAction = "update_config"
	ActionStart        TaskAction = "start"
	ActionStop         TaskAction = "stop"
	ActionReloadConfig TaskAction = "reload_config"
)

type TaskStatus string
type TaskAction string
type TaskType string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

const (
	TaskTypeGo         TaskType = "go"
	TaskTypePrometheus TaskType = "prometheus"
	TaskTypeShell      TaskType = "shell"
)

type RunnerConfig struct {
	RunnerID       string   `yaml:"runner_id"`
	ManagerURL     string   `yaml:"manager_url"`
	Applications   []string `yaml:"applications"`
	Groups         []string `yaml:"groups"`
	SupportedTasks []string `yaml:"supported_tasks"`
	RunnerAddress  string   `yaml:"runner_address"`
	ListenAddr     string   `yaml:"addr"`
}

type Task struct {
	ID     string     `json:"id"`
	Action TaskAction `json:"action"`
	Type   TaskType   `json:"type"`
	Target string     `json:"target"`
	Config struct {
		Content string `json:"content"`
		Path    string `json:"path"`
	} `json:"config,omitempty"`
	Script      string     `json:"script,omitempty"`
	ProcessArgs []string   `json:"process_args,omitempty"`
	Status      TaskStatus `json:"status"`
	Error       string     `json:"error,omitempty"`
	Dir         string     `json:"dir,omitempty"`
}

type TaskStatusUpdate struct {
	TaskID    string     `json:"task_id"`
	Status    TaskStatus `json:"status"`
	Error     string     `json:"error,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

type TaskProcessor interface {
	Submit(task Task) error
	GetStatus(taskID string) (TaskStatus, error)
}
type StatusUpdater interface {
	UpdateStatus(taskID string, status TaskStatus, errMsg string) error
}
type TaskHandler interface {
	Execute(task *Task) error
}
