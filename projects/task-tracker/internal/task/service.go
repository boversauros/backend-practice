package task

import (
	"fmt"
	"time"
)

type Service struct {
    storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetTasks(status *Status) ([]Task, error) {
    tasks, err := s.storage.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading tasks: %v", err)
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

    return matchingTasks, nil
}



func (s *Service) AddTask(description string) (int, error){
	tasks, err := s.storage.ReadAll()
    if err != nil {
        return 0, fmt.Errorf("error reading tasks: %v", err)
    }

	now := time.Now()
	newTask := Task{
		ID: len(tasks) + 1,
		Description: description,
		Status: StatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}
    
	tasks = append(tasks, newTask)
	err = s.storage.Save(tasks)
	if err != nil {
		return 0, fmt.Errorf("error saving tasks: %v", err)
	}

	return newTask.ID, nil
}

func (s *Service) UpdateTaskDescription(id int, description string) error {
    tasks, err := s.storage.ReadAll()
    if err != nil {
        return fmt.Errorf("error reading tasks: %v", err)
    }

    taskFound := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
            tasks[i].UpdatedAt = time.Now()
            taskFound = true
            break
		}
	}

    if taskFound {
        err = s.storage.Save(tasks)
        if err != nil {
            return fmt.Errorf("error saving tasks: %v", err)
        }
    } else {
		return fmt.Errorf("Error: Task with ID %d not found\n", id)
    }

	return nil
}

func (s *Service) UpdateTaskStatus(id int, status Status) error {
    tasks, err := s.storage.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading tasks: %v", err)
	}
        
    taskFound := false
	for i, task := range tasks {
		if task.ID == id {
            tasks[i].Status = status
            tasks[i].UpdatedAt = time.Now()
            taskFound = true
            break
		}
	}

    if taskFound {
        err = s.storage.Save(tasks)
		if err != nil {
			return fmt.Errorf("error saving tasks: %v", err)
		}
    } else {
		return fmt.Errorf("Error: Task with ID %d not found\n", id)
    }

	return nil
}


func (s *Service) DeleteTask(id int) error {
    tasks, err := s.storage.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading tasks: %v", err)
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
        err = s.storage.Save(tasks)
        if err != nil {
            return fmt.Errorf("error saving tasks: %v", err)
        }
    } else {
        return fmt.Errorf("Error: Task with ID %d not found\n", id)
    }

	return nil
}
