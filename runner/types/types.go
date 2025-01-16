// Config represents the runner configuration
package types

type RunnerConfig struct {
	RunnerID       string   `yaml:"runner_id"`
	ManagerURL     string   `yaml:"manager_url"`
	Applications   []string `yaml:"applications"`
	Groups         []string `yaml:"groups"`
	SupportedTasks []string `yaml:"supported_tasks"`
	RunnerAddress  string   `yaml:"runner_address"`
}
