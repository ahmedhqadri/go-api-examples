package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
)

func AddTask(respWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var task Task
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(respWriter, "Show Details for ToDo '%v'\n", task)

	accessTasks.Lock()
	task.ID = allTasks[len(allTasks)-1].ID + 1

	allTasks = append(allTasks, task)
	accessTasks.Unlock()

	saveCSV()
}

type Check struct {
	ID int `json:"id"`
}

func CheckOff(respWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var check Check
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, &check)
	if err != nil {
		log.Println(err)
		return
	}

	checkTaskAsComplete(check.ID)
}

func checkTaskAsComplete(taskID int) (err error) {
	// Get the index.
	index, err := getIndexByTaskID(taskID)
	if err != nil {
		log.Println(err)
		return
	}

	// Lock the tasks.
	accessTasks.Lock()

	// Set the task to checked.
	allTasks[index].Checked = true

	accessTasks.Unlock()
	return
}
