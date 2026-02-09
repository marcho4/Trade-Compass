package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := Load("../../tasks.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(cfg.Tasks))
	}

	task1 := cfg.Tasks[0]
	if task1.Name != "find_new_reports" {
		t.Errorf("Expected task name 'find_new_reports', got '%s'", task1.Name)
	}
	if task1.Cron != "0 9 * * *" {
		t.Errorf("Expected cron '0 9 * * *', got '%s'", task1.Cron)
	}
	if task1.Topic != "parser.find_reports" {
		t.Errorf("Expected topic 'parser.find_reports', got '%s'", task1.Topic)
	}
	if task1.Message["action"] != "find_new_reports" {
		t.Errorf("Expected action 'find_new_reports', got '%v'", task1.Message["action"])
	}
}
