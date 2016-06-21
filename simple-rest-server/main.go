package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Task struct {
	ID        int
	Checked   bool      `json:"checked"`
	TimeAdded time.Time `json:"time_added"`
	Deadline  time.Time `json:"deadline"`
	Task      string    `json:"task"`
}

var (
	allTasks    []Task
	accessTasks = &sync.Mutex{}
)

const (
	timeFormat = time.RFC3339
)

func main() {

	// Loading the csv file into the RAM
	csvfile, err := os.Open("tasks.csv")
	ifPanic(err)
	rawCSVdata, err := csv.NewReader(csvfile).ReadAll()
	ifPanic(err)

	for i, each := range rawCSVdata {
		timeAdded, err := time.Parse(timeFormat, each[2])
		ifPanic(err)
		deadline, err := time.Parse(timeFormat, each[3])
		ifPanic(err)
		status := false
		if each[1] == "true" {
			status = true
		}
		allTasks = append(allTasks, Task{ID: i, Task: each[4], TimeAdded: timeAdded, Deadline: deadline, Checked: status})
	}

	// Autosave every minutes
	csvfile.Close()
	go func() {
		for {
			saveCSV()
			time.Sleep(1 * time.Minute)
		}
	}()

	// Start the API
	router := httprouter.New()
	router.GET("/search", SearchTask)
	router.GET("/list", ListTask)
	router.POST("/add", AddTask)
	router.DELETE("/delete", DeleteTask)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func saveCSV() {
	myString := ""
	accessTasks.Lock()
	for _, each := range allTasks {
		myString += fmt.Sprintf("%v,%v,%v,%v,\"%v\"\n", each.ID, each.Checked, each.TimeAdded.Format(timeFormat), each.Deadline.Format(timeFormat), each.Task)
	}
	ioutil.WriteFile("tasks.csv", []byte(myString), 0644)
	accessTasks.Unlock()
}

func ifPanic(err error) {
	if err != nil {
		// Maybe change to log later
		fmt.Println(err)
		os.Exit(1)
	}
}
