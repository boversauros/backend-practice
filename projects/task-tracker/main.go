package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TaskStatus string

const (
    Todo TaskStatus = "todo"
    InProgress TaskStatus = "in_progress"
    Done TaskStatus = "done"
)


type Task struct {
	id 		int
	Description string
	Status TaskStatus  
}

var tasks []Task

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    
    fmt.Println("Welcome to Task CLI! Type 'help' to see available commands.")
    fmt.Print("task-cli > ")
    
    for scanner.Scan() {
        userInput := scanner.Text()
        words := strings.Fields(userInput)
        
        if len(words) == 0 {
            fmt.Print("task-cli > ")
            break
        }
        
        command := words[0]
        parameters := words[1:]

        switch command {
        case "list", "ls":
            if len(tasks) == 0 {
                fmt.Println("You have no tasks!")
            } else { 
                listTasks()
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
            addTask(taskDescription)
            
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
            taskDescription := strings.Join(parameters[1:], " ")
            updateTaskDescription(taskID, taskDescription)

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
            updateTaskStatus(taskID, InProgress)

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
            updateTaskStatus(taskID, Done)
        
        case "delete":
            if len(parameters) == 0 {
                fmt.Println("Error: Please provide a task ID")
                break
            } 
            taskID, err := strconv.Atoi(parameters[0])
            if err != nil {
                fmt.Println("Error: Invalid task ID")
                break   
            }
            deleteTask(taskID)
        
        case "help":
            fmt.Println("Available commands:")
            fmt.Println("  list or ls - List all tasks")
            fmt.Println("  add <task description> - Add a new task")
            fmt.Println("  update <task ID> <new task description> - Update a task description")
            fmt.Println("  mark-in-progress <task ID> - Mark a task as in progress")
            fmt.Println("  mark-done <task ID> - Mark a task as done")
            fmt.Println("  delete <task ID> - Delete a task")            
            fmt.Println("  help - Show this help message")
            fmt.Println("  exit - Exit the program")
            
        case "exit":
            fmt.Println("Goodbye!")
            return
            
        default:
            fmt.Printf("Unknown command: %q\n", command)
            fmt.Println("Type 'help' to see available commands")
        }
        
        fmt.Print("task-cli > ")
    }
}


func addTask(description string) {
	id := len(tasks) + 1
	tasks = append(tasks, Task{id: id, Description: description, Status: "todo"})
    fmt.Printf("Task added successfully (ID: %d)\n", id)
}

func updateTaskDescription(id int, description string) {
    taskFound := false

	for i, task := range tasks {
		if task.id == id {
			tasks[i] = Task{id: id, Description: description}
            taskFound = true
            break
		}
	}

    if !taskFound {
        fmt.Printf("Error: Task with ID %d not found\n", id)
    }
}

func updateTaskStatus(id int, status TaskStatus) {
    taskFound := false

	for i, task := range tasks {
		if task.id == id {
			tasks[i] = Task{id: id, Status: status}
            taskFound = true
            break
		}
	}

    if !taskFound {
        fmt.Printf("Error: Task with ID %d not found\n", id)
    }
}

func deleteTask(id int) {
    taskFound := false

	for i, task := range tasks {
		if task.id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
            taskFound = true
            break
		} 
	}

    if !taskFound {
        fmt.Printf("Error: Task with ID %d not found\n", id)
    }
}

func listTasks() {
	for _, task := range(tasks) {
		fmt.Printf("%d: %s - %s\n", task.id, task.Description, task.Status)
	}
}