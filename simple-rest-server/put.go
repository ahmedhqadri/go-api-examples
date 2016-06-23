package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Check struct {
	ID int `json:"id"`
}

func Modify(respWriter http.ResponseWriter, request *http.Request, ps httprouter.Params) {
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
	ID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		// TODO: display 400, id is not a number
		return
	}

	// TODO: be able to change anything
	checkTaskAsComplete(ID)
	// TODO: error 400, print it on the body
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
