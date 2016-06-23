package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
)

type Remove struct {
	ID int `json:"id"`
}

// Delete API call
func DeleteTask(respWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var remove Remove
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintln(respWriter, "Bad request: %v", err)
		log.Println(err)
		respWriter.WriteHeader(400)
	}
	err = json.Unmarshal(body, &remove)
	if err != nil {
		fmt.Fprintln(respWriter, "Bad request: %v", err)
		log.Println(err)
		respWriter.WriteHeader(400)
	}

	index, err := removeTaskByID(remove.ID)
	if index == -1 {
		fmt.Fprintln(respWriter, "Task ID not found. Unable to delete")
	} else {
		fmt.Fprintln(respWriter, "Task ID: "+string(index)+" found.")
	}
	if err != nil {
		fmt.Fprintln(respWriter, "Bad request: %v", err)
		log.Println(err)
		respWriter.WriteHeader(400)
	}
}

func removeTaskByID(taskID int) (int, error) {
	index, err := getIndexByTaskID(taskID)
	if err != nil {
		log.Println(err)
		return index, err
	}

	if index == -1 {
		return index, err
	} else if index == 0 {
		accessTasks.Lock()
		allTasks = allTasks[1:]
		accessTasks.Unlock()
	} else if index == len(allTasks)-1 {
		accessTasks.Lock()
		allTasks = allTasks[:len(allTasks)-1]
		accessTasks.Unlock()
	} else {
		accessTasks.Lock()
		allTasks = append(allTasks[:index-1], allTasks[index+1:]...)
		accessTasks.Unlock()
	}
	accessTasks.Unlock()
	return index, err
}
