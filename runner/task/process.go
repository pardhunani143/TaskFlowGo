package task

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pardhunani143/TaskFlowGo/runner/types"
)

type ProcessConfig struct {
	Path       string   `json:"path"`
	WorkingDir string   `json:"working_dir"`
	Args       []string `json:"args"`
}

func writeConfig(filename, content string) error {
	tempFile := filename + ".tmp"
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return err
	}
	return os.Rename(tempFile, filename)
}

func startProcess(config ProcessConfig) error {
	if config.Path == "" {
		return fmt.Errorf("process path is required")
	}
	log.Println("Starting process", config.Path)
	log.Println("With args,", config.Args)
	cmd := exec.Command(config.Path, config.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("cmd.Stdout", cmd.Stdout)
	log.Println("cmd.Stderr", cmd.Stderr)

	if config.WorkingDir != "" {
		cmd.Dir = config.WorkingDir
	}

	return cmd.Start()
}

func stopProcess(processName string) error {
	// Check if process exists first
	checkCmd := exec.Command("pgrep", "-f", processName)
	if err := checkCmd.Run(); err != nil {
		// Process not found - return nil as it's already stopped
		return nil
	}

	// Process exists, try to stop it
	stopCmd := exec.Command("pkill", "-f", processName)
	return stopCmd.Run()
}

func reloadProcess(processName string) error {
	cmd := exec.Command("pkill", "-HUP", processName)
	return cmd.Run()
}

type GoHandler struct{}

func (h *GoHandler) Execute(task *types.Task) error {
	if task.Config.Content != "" {
		if err := writeConfig(task.Config.Path, task.Config.Content); err != nil {
			return err
		}
	}

	switch task.Action {
	case "start", "restart":
		if err := stopProcess("go"); err != nil {
			return err
		}
		args := append([]string{"--config", task.Config.Path}, task.ProcessArgs...)
		procConfig := ProcessConfig{
			Path: "./process",
			Args: args,
		}
		return startProcess(procConfig)
	case "stop":
		return stopProcess("go")
	}
	return nil
}

type PrometheusHandler struct{}

func (h *PrometheusHandler) Execute(task *types.Task) error {
	switch task.Action {
	case "start":
		args := append([]string{"--config.file", task.Config.Path}, task.ProcessArgs...)
		procConfig := ProcessConfig{
			Path:       "./prometheus",
			Args:       args,
			WorkingDir: task.Dir,
		}
		return startProcess(procConfig)

	case "stop":
		return stopProcess("prometheus")

	case "restart":
		if err := stopProcess("prometheus"); err != nil {
			return err
		}

		args := append([]string{"--config.file", task.Config.Path}, task.ProcessArgs...)
		procConfig := ProcessConfig{
			Path:       "./prometheus",
			Args:       args,
			WorkingDir: task.Dir,
		}
		return startProcess(procConfig)

	case "update_config":
		if err := writeConfig(task.Config.Path, task.Config.Content); err != nil {
			return err
		}
		return reloadProcess("prometheus")

	case "reload":
		return reloadProcess("prometheus")

	default:
		return fmt.Errorf("unknown action: %s", task.Action)
	}
}

type ShellHandler struct{}

func (h *ShellHandler) Execute(task *types.Task) error {
	if task.Script == "" {
		return nil
	}

	scriptFile := "/tmp/script_" + task.ID + ".sh"
	if err := writeConfig(scriptFile, task.Script); err != nil {
		return err
	}
	defer os.Remove(scriptFile)

	if err := os.Chmod(scriptFile, 0755); err != nil {
		return err
	}

	cmd := exec.Command(scriptFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
