package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	Name  string
	Detail string
}

var taskList = []Task{
	{Name: "Task 1", Detail: "Detail 1"},
	{Name: "Task 2", Detail: "Detail 2"},
	{Name: "Task 3", Detail: "Detail 3"},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("tmpl").Parse(`<ul class='list-group'>{{range .}}<li class='list-group-item bg-primary text-white'>{{.Name}} - {{.Detail}}</li>{{end}}</ul>`)
	tmpl.Execute(w, taskList)
}

func main() {
	fmt.Println("Go app...")
    
	// get tasks
	h1 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
			return
		}

		getTasks(w, r)
	}

	// get task
	h2 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
			return
		}

		flag := false
		name := r.URL.Query().Get("name")
		for _, task := range taskList {
			if task.Name == name {
				htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", task.Name, task.Detail)
				tmpl, _ := template.New("tmpl").Parse(htmlStr)
				tmpl.Execute(w, task)
				flag = true
				break
			}
		}

		if !flag {
			http.Error(w, "Task not found!", http.StatusNotFound)
		}
	}

	// add new task
	h3 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			newTask := Task{Name: r.FormValue("name"), Detail: r.FormValue("detail")}
			taskList = append(taskList, newTask)
			htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", newTask.Name, newTask.Detail)
			tmpl, _ := template.New("tmpl").Parse(htmlStr)
			tmpl.Execute(w, newTask)
		}else {
			http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		}
	}

	// delete task
	h4 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
			return
		}

		flag := false
		name := r.URL.Query().Get("name")
		for index, task := range taskList {
			if task.Name == name {
				taskList = append(taskList[:index], taskList[index+1:]...)
				flag = true
				break
			}
		}

		if !flag {
			http.Error(w, "Task not found!", http.StatusNotFound)
		}else {
			// get tasks
			getTasks(w, r)
		}
	}

	// define handlers
	http.HandleFunc("/getTasks", h1)
	http.HandleFunc("/getTask", h2)
	http.HandleFunc("/createTask", h3)
	http.HandleFunc("/deleteTask", h4)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
