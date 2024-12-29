package main

import (
	"os"
	"task-tracker/internal/cli"
	"task-tracker/internal/task"
)

func main() {
    // Create a directory to store our tasks
    if err := os.MkdirAll("data", 0755); err != nil {
        panic(err)
    }
    
    // Create a file storage for our tasks
    tasksFile := "data/tasks.json"
    storage := task.NewFileStorage(tasksFile)

    // Create our service with the storage
    service := task.NewService(storage)

    // Create and run our CLI handler
    handler := cli.NewHandler(service)
    handler.Run()
}