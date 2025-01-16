package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pardhunani143/TaskFlowGo/runner/config"
	"github.com/pardhunani143/TaskFlowGo/runner/types"
)

func main() {
	config, err := config.LoadConfig("config.yml")

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Register with Manager
	err = registerRunner(config)
	if err != nil {
		log.Fatalf("Failed to register Runner: %v", err)
	}

	// Start Heartbeat
	go startHeartbeat(config)

	// Wait for tasks
	log.Println("Runner is ready and waiting for tasks...")
	select {} // Keep the program running
}

func registerRunner(config types.RunnerConfig) error {
	payload := types.RunnerConfig{
		RunnerID:       config.RunnerID,
		Applications:   config.Applications,
		Groups:         config.Groups,
		SupportedTasks: config.SupportedTasks,
		RunnerAddress:  config.RunnerAddress,
	}

	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(config.ManagerURL+"/register", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	log.Println("Runner registered successfully.")
	return nil
}

func startHeartbeat(config types.RunnerConfig) {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		payload := map[string]string{
			"runner_id": config.RunnerID,
			"status":    "healthy",
		}
		payloadBytes, _ := json.Marshal(payload)

		resp, err := http.Post(config.ManagerURL+"/heartbeat", "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to send heartbeat: %v", err)
		} else {
			log.Println("Heartbeat sent successfully.")
		}
	}
}
