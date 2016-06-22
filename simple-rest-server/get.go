package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Search struct {
	Word       string    `json:"search"`
	TimeBefore time.Time `json:"before"`
}

func Index(respWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	htmlOut := `Todo List API.

<p>/add</p>
<p>/delete</p>
<p>/search</p>
<p>/list</p>
`

	respWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	respWriter.Write([]byte(htmlOut))
}

func ListTask(respWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	output, _ := json.MarshalIndent(allTasks, "", "  ")
	fmt.Fprintln(respWriter, string(output))
}

func SearchTask(respWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var query Search
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &query)
	if err != nil {
		fmt.Println(err)
	}

	var indicies []int
	timezero := time.Time{}
	if query.Word != "" {
		indicies = append(indicies, searchByWord(query.Word)[:]...)
	} else if query.TimeBefore != timezero {
		indicies = append(indicies, searchByTime(query.TimeBefore)[:]...)
	} else {
		fmt.Println("Undefined search query")
	}
}

func searchByWord(query string) (result []int) {
	accessTasks.Lock()
	for i, task := range allTasks {
		if strings.Contains(task.Task, query) {
			result = append(result, i)
		}
	}
	accessTasks.Unlock()
	return
}

func searchByTime(query time.Time) (result []int) {
	accessTasks.Lock()
	for i, task := range allTasks {
		if task.TimeAdded.Before(query) || task.TimeAdded == query {
			result = append(result, i)
		}
	}
	accessTasks.Unlock()
	return
}
