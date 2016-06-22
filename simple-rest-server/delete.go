package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Remove struct {
	ID int `json:"id"`
}

// Delete API call
func DeleteTask(respWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var remove Remove
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &remove)
	if err != nil {
		log.Println(err)
	}

	removeTaskByID(remove.ID)
	// TODO: error 400, print it on the body
}

func removeTaskByID(taskID int) (err error) {
	index, err := getIndexByTaskID(taskID)
	if err != nil {
		log.Println(err)
		return
	}
	accessTasks.Lock()
	allTasks = append(allTasks[:index-1], allTasks[index+1:]...)
	accessTasks.Unlock()
	return
}
