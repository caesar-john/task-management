package handler

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/caesar-john/task-management/entities"
	"github.com/caesar-john/task-management/repository"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/partials/header.html",
		"templates/partials/sidebar.html",
		"templates/pages/dashboard.html",
	)

	if err != nil {
		http.Error(w, "Failed to received the data", http.StatusInternalServerError)
		return
	}

	tasks, err := repository.GetAllTask(r.Context())
	if err != nil {
		http.Error(w, "Failed to received the data", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	greeting := "Good evening!"
	if now.Hour() < 12 {
		greeting = "Good morning!"
	} else if now.Hour() < 18 {
		greeting = "Good afternoon!"
	}

	viewData := entities.DashboardViewData{
		Tasks:      tasks,
		TotalTasks: len(tasks),
		Filter:     r.URL.Query().Get("status"),
		DateLabel:  strings.ToUpper(now.Format("Monday, January 2, 2006")),
		Greeting:   greeting,
	}
	if viewData.Filter != "Pending" && viewData.Filter != "Completed" {
		viewData.Filter = "All"
	}

	for _, task := range tasks {
		if strings.EqualFold(task.Status, "Completed") {
			viewData.CompletedTasks++
		} else {
			viewData.PendingTasks++
		}
	}

	if viewData.Filter != "All" {
		filteredTasks := make([]entities.Task, 0)
		for _, task := range viewData.Tasks {
			if strings.EqualFold(task.Status, viewData.Filter) {
				filteredTasks = append(filteredTasks, task)
			}
		}
		viewData.Tasks = filteredTasks
	}

	err = tpl.ExecuteTemplate(w, "base", viewData)
	if err != nil {
		http.Error(w, "Failed to received the data", http.StatusInternalServerError)
		return
	}
}
