package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
)

type Data struct {
	TaskIdCounter int          `json:"task_id_counter"`
	Tasks         map[int]Task `json:"tasks"`
}

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

var data = Data{0, map[int]Task{}}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	loadStoredData()
	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addCmd.Parse(os.Args[2:])
		positionalArgs := addCmd.Args()
		if len(positionalArgs) != 1 {
			os.Exit(2)
		}
		addTask(positionalArgs[0])

	case "list":
		list()

	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		updateCmd.Parse(os.Args[2:])
		positionalArgs := updateCmd.Args()
		if len(positionalArgs) < 2 {
			os.Exit(2)
		}
		id, err := strconv.Atoi(positionalArgs[0])
		if err != nil {
			os.Exit(1)
		}
		updateTask(id, positionalArgs[1])

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteCmd.Parse(os.Args[2:])
		positionalArgs := deleteCmd.Args()
		if len(positionalArgs) < 1 {
			os.Exit(2)
		}
		id, err := strconv.Atoi(positionalArgs[0])
		if err != nil {
			os.Exit(1)
		}
		deleteTask(id)
	}

}

// Data Store Methods

func createDataStoreFileIfNeeded() {
	if _, error := os.Stat("data.json"); errors.Is(error, os.ErrNotExist) {
		file, create_file_error := os.Create("data.json")
		if create_file_error != nil {
			log.Println("[error] ", create_file_error)
			os.Exit(1)
		}
		file.WriteString("{}")
		file.Close()
	}
}

func loadStoredData() {
	createDataStoreFileIfNeeded()

	storedDataBytes, error := os.ReadFile("data.json")
	if error != nil {
		log.Println("unable to to parse stored data - ", error)
		return
	}

	if storedDataBytes != nil {
		if error := json.Unmarshal(storedDataBytes, &data); error != nil {
			log.Println("failed to unmarshal stored data - ", error)
		}
	}
}

func save() {
	jsonBytes, marshal_error := json.Marshal(data)
	if marshal_error != nil {
		log.Println("failed to marshal data - ", marshal_error)
	}
	file, open_file_error := os.OpenFile("data.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if open_file_error != nil {
		log.Println("failed to  open data file for save - ", open_file_error)
	}
	defer file.Close()

	_, write_error := file.WriteString(string(jsonBytes))
	if write_error != nil {
		log.Println("failed to write json to data file - ", write_error)
	}
}

func getNewTaskId() int {
	id := data.TaskIdCounter
	data.TaskIdCounter += 1
	return id
}

// Actions

func addTask(description string) {
	id := getNewTaskId()
	data.Tasks[id] = Task{id, description}
	save()
}

func list() {
	tasksSlice := []Task{}
	for _, task := range data.Tasks {
		tasksSlice = append(tasksSlice, task)
	}
	bytes, _ := json.Marshal(tasksSlice)
	print(string(bytes))
}

func updateTask(id int, description string) {
	if task, ok := data.Tasks[id]; ok {
		task.Description = description
		data.Tasks[id] = task
		save()
	}
}

func deleteTask(id int) {
	delete(data.Tasks, id)
	save()
}
