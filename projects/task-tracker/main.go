package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TaskStatus string

const (
    StatusTodo TaskStatus = "todo"
    StatusInProgress TaskStatus = "in_progress"
    StatusDone TaskStatus = "done"
)

const tasksFile = "tasks.json"


type Task struct {
    ID 		int  `json:"id"`
	Description string `json:"description"`
	Status TaskStatus  `json:"status"`
}

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
            tasks, err := readTasks()
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                break
            }

            if len (tasks) == 0 {
                fmt.Println("No tasks found")
                break
            }

            if len(parameters) == 0 {
                listTasks(nil)
            } else {
                status := TaskStatus(parameters[0])
                listTasks(&status)
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
            updateTaskStatus(taskID, StatusInProgress)

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
            updateTaskStatus(taskID, StatusDone)
        
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
            return
            
        default:
            fmt.Printf("Unknown command: %q\n", command)
            fmt.Println("Type 'help' to see available commands")
        }
        
        fmt.Print("task-cli > ")
    }
}

func readTasks() ([]Task, error) {
    if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
        return []Task{}, nil
    }

    data, err := os.ReadFile(tasksFile)
    if err != nil {
        return nil, err
    }
    
    if len(data) == 0 {
        return []Task{}, nil
    }    
    
    var tasks []Task
	err = json.Unmarshal(data, &tasks)
    if err != nil {
        return nil, fmt.Errorf("error reading tasks: %v", err)
    }

    return tasks, nil
}

func writeTasks(tasks []Task) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return fmt.Errorf("error reading tasks: %v", err)
    }

    err = os.WriteFile(tasksFile, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing tasks: %v", err)
    }

    return nil
}

func addTask(description string) {
    tasks, err := readTasks()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

	id := len(tasks) + 1
	tasks = append(tasks, Task{ID: id, Description: description, Status: StatusTodo})
    err = writeTasks(tasks)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Task added successfully (ID: %d)\n", id)
}

func updateTaskDescription(id int, description string) {
    tasks, err := readTasks()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    taskFound := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = Task{ID: id, Description: description}
            taskFound = true
            break
		}
	}

    if taskFound {
        err = writeTasks(tasks)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
    } else {
        fmt.Printf("Error: Task with ID %d not found\n", id)
    }
}

func updateTaskStatus(id int, status TaskStatus) {
    tasks, err := readTasks()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    
    taskFound := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = Task{ID: id, Status: status}
            taskFound = true
            break
		}
	}

    if taskFound {
        err = writeTasks(tasks)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
    } else {
        fmt.Printf("Error: Task with ID %d not found\n", id)

    }
}

func deleteTask(id int) {
    tasks, err := readTasks()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    taskFound := false
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
            taskFound = true
            break
		} 
	}

    if taskFound {
        err = writeTasks(tasks)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
    } else {
        fmt.Printf("Error: Task with ID %d not found\n", id)
    }
}

func listTasks(status *TaskStatus) {
    tasks, err := readTasks()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    var matchingTasks []Task

    if status == nil {
        matchingTasks = tasks
    } else {
        for _, task := range tasks {
            if task.Status == *status {
                matchingTasks = append(matchingTasks, task)
            }
        }
    }

    if len(matchingTasks) == 0 {
        fmt.Println("No tasks found")
        return
    }

    for _, task := range matchingTasks {
        fmt.Printf("%d: %s [%s]\n", task.ID, task.Description, task.Status)
    }
}