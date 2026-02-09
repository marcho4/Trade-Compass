package main

import (
	"log"

	"schedulerservice/internal/config"
)

func main() {
	cfg, err := config.Load("tasks.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded %d tasks from config", len(cfg.Tasks))
	for _, task := range cfg.Tasks {
		log.Printf("- Task: %s, Cron: %s, Topic: %s", task.Name, task.Cron, task.Topic)
	}
}
