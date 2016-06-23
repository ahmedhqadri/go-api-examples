package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddTask(respWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var newTask Task
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &newTask)
	if err != nil {
		fmt.Fprintf(respWriter, "Bad request: %v", err)
		respWriter.WriteHeader(400)
		return
	}

	accessTasks.Lock()
	newTask.ID = 0
	if len(allTasks) > 0 {
		newTask.ID = allTasks[len(allTasks)-1].ID + 1
	}
	// lastID := allTasks[len(allTasks)-1].ID
	// for i := range newTask {
	// 	lastID++
	// 	newTask[i].ID = lastID + 1
	// }

	// allTasks = append(allTasks, newTask...)
	newTask.ID = allTasks[len(allTasks)-1].ID + 1
	allTasks = append(allTasks, newTask)
	accessTasks.Unlock()

	err = saveCSV()
	if err != nil {
		fmt.Fprintf(respWriter, "Internal error: %v", err)
		respWriter.WriteHeader(500)
		return
	}
}
