package task

import (
	"encoding/json"
	"fmt"
	"os"
)


type Storage interface {
	ReadAll() ([]Task, error)
	Save(tasks []Task) error
}

type FileStorage struct {
	filePath string
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{filePath: filePath}
}

func (fs *FileStorage) ReadAll() ([]Task, error) {
	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
        return []Task{}, nil
    }

    data, err := os.ReadFile(fs.filePath)
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

func (fs *FileStorage) Save(tasks []Task) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return fmt.Errorf("error reading tasks: %v", err)
    }

    err = os.WriteFile(fs.filePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing tasks: %v", err)
    }

    return nil
}