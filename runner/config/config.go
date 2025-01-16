package config

import (
	"log"
	"os"

	"github.com/pardhunani143/TaskFlowGo/runner/types"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from a file
func LoadConfig(filepath string) (types.RunnerConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var config types.RunnerConfig
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return types.RunnerConfig{}, err
	}

	return config, nil
}
