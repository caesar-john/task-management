package handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/caesar-john/task-management/entities"
	"github.com/caesar-john/task-management/repository"
)

type taskFormData struct {
	DateLabel   string
	Greeting    string
	Filter      string
	Title       string
	Description string
	Deadline    string
	Error       string
}

func GetTaskForm(w http.ResponseWriter, r *http.Request) {
	data := newTaskFormData()
	renderTaskForm(w, data, http.StatusOK)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		data := newTaskFormData()
		data.Error = "Please check the form and try again."
		renderTaskForm(w, data, http.StatusBadRequest)
		return
	}

	data := newTaskFormData()
	data.Title = strings.TrimSpace(r.FormValue("title"))
	data.Description = strings.TrimSpace(r.FormValue("description"))
	data.Deadline = r.FormValue("deadline")
	if data.Title == "" || data.Deadline == "" {
		data.Error = "A title and deadline are required."
		renderTaskForm(w, data, http.StatusBadRequest)
		return
	}

	deadline, err := time.Parse("2006-01-02", data.Deadline)
	if err != nil {
		data.Error = "Enter a valid deadline."
		renderTaskForm(w, data, http.StatusBadRequest)
		return
	}

	now := time.Now()
	task := entities.Task{
		Title:       data.Title,
		Description: data.Description,
		Status:      "Pending",
		Deadline:    deadline,
		CreatedAt:   now,
	}
	if err := repository.CreateTask(r.Context(), task); err != nil {
		log.Printf("Failed to create task: %v", err)
		data.Error = "Could not save the task. Please try again."
		renderTaskForm(w, data, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func newTaskFormData() taskFormData {
	now := time.Now()
	greeting := "Good evening!"
	if now.Hour() < 12 {
		greeting = "Good morning!"
	} else if now.Hour() < 18 {
		greeting = "Good afternoon!"
	}

	return taskFormData{
		DateLabel: strings.ToUpper(now.Format("Monday, January 2, 2006")),
		Greeting:  greeting,
		Filter:    "All",
	}
}

func renderTaskForm(w http.ResponseWriter, data taskFormData, status int) {
	tpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/partials/header.html",
		"templates/partials/sidebar.html",
		"templates/pages/task.html",
	)
	if err != nil {
		log.Printf("Failed to load task form: %v", err)
		http.Error(w, "Could not load the task form.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if err := tpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Failed to render task form: %v", err)
	}
}
