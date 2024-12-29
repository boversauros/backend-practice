package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"task-tracker/internal/task"
)


type Handler struct {
	service *task.Service
	scanner *bufio.Scanner
}

func NewHandler(service *task.Service) *Handler {
	return &Handler{
		service: service, 
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (h *Handler) Run() {
	fmt.Println("Welcome to Task CLI! Type 'help' to see available commands.")

    for {
        fmt.Print("task-cli > ")
        if !h.scanner.Scan() {
            break
        }

        if err := h.handleCommand(h.scanner.Text()); err != nil {
            fmt.Printf("Error: %v\n", err)
		}
	}
}

func (h *Handler) handleCommand(userInput string) error {
	words := strings.Fields(userInput)

	if len(words) == 0 {
		return nil
	}

	command := words[0]
	parameters := words[1:]

	switch command {
	case "list", "ls":
		var status *task.Status
		if len(parameters) > 0 {
			s := task.Status(parameters[0])
			status = &s
		}

		tasks, err := h.service.GetTasks(status)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			break
		}

		for _, t := range tasks {
			fmt.Printf("%d. %s [%s]\n", t.ID, t.Description, t.Status)
		}

	case "add":
		if len(parameters) == 0 {
			fmt.Println("Error: Please provide a task description")
			break
		}

		input := strings.Join(parameters, " ")
		re := regexp.MustCompile(`"([^"]*)"`)
		matches := re.FindStringSubmatch(input)

		if len(matches) < 2 {
			fmt.Println("Error: Task description must be enclosed in double quotes")
			break
		}

		taskDescription := matches[1] 
		if taskDescription == "" { 
			fmt.Println("Error: Task description cannot be empty")
			break
		}
		id, err := h.service.AddTask(taskDescription)
		if err != nil {
			return fmt.Errorf("error adding task: %v", err)
		}

		fmt.Printf("Task added successfully (ID: %d)\n", id)
	case "update":
		if len(parameters) < 2 {
			fmt.Println("Error: Please provide a task ID and description")
			break
		}  

		taskID, err := strconv.Atoi(parameters[0])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			break
		}
		input := strings.Join(parameters[1:], " ")
		re := regexp.MustCompile(`"([^"]*)"`)
		matches := re.FindStringSubmatch(input)

		if len(matches) < 2 {
			fmt.Println("Error: Task description must be enclosed in double quotes")
			break
		}

		taskDescription := matches[1]
		if taskDescription == "" {
			fmt.Println("Error: Task description cannot be empty")
			break
		}

		err = h.service.UpdateTaskDescription(taskID, taskDescription)
		if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}
	
	case "mark-in-progress":
		if len(parameters) == 0 {
			fmt.Println("Error: Please provide a task ID")
			break
		}
		taskID, err := strconv.Atoi(parameters[0])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			break
		}

		err = h.service.UpdateTaskStatus(taskID, task.StatusInProgress)
		if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}
	
	case "mark-done":
		if len(parameters) == 0 {
			fmt.Println("Error: Please provide a task ID")
			break
		}
		taskID, err := strconv.Atoi(parameters[0])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			break
		}
		err = h.service.UpdateTaskStatus(taskID, task.StatusDone)
		if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}

	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  list - List all tasks")
		fmt.Println("  list <status> - List tasks with a specific status")
		fmt.Println("  add <task description> - Add a new task")
		fmt.Println("  update <task ID> <new task description> - Update a task description")
		fmt.Println("  mark-in-progress <task ID> - Mark a task as in progress")
		fmt.Println("  mark-done <task ID> - Mark a task as done")
		fmt.Println("  delete <task ID> - Delete a task")            
		fmt.Println("  help - Show this help message")
		fmt.Println("  exit - Exit the program")
		
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
		
	default:
		fmt.Printf("Unknown command: %q\n", command)
		fmt.Println("Type 'help' to see available commands")
	}

	return nil
}