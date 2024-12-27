package main

import (
	"path/filepath"
	"runtime"
	"task-tracker/internal/cli"
	"task-tracker/internal/task"
)

func main() {
	// Get the path to the tasks.json file
    _, currentFile, _, _ := runtime.Caller(0)
    sourceDir := filepath.Dir(currentFile)
    tasksFile := filepath.Join(sourceDir, "tasks.json")
    storage := task.NewFileStorage(tasksFile)

    // Create our service with the storage
    service := task.NewService(storage)

    // Create and run our CLI handler
    handler := cli.NewHandler(service)
    handler.Run()
}