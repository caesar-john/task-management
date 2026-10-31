package main

import (
	"log"
	"net/http"
	"os"

	"github.com/caesar-john/task-management/config"
	"github.com/caesar-john/task-management/handler"
)

func main() {
	config.ConnectDb()

	//static
	fs := http.FileServer(http.Dir("statics"))

	//dasboard
	http.HandleFunc("GET /dashboard", handler.GetDashboard)
	http.HandleFunc("GET /task", handler.GetTaskForm)
	http.HandleFunc("POST /task", handler.CreateTask)
	http.HandleFunc("POST /task/{id}/toggle", handler.ToggleTaskStatus)
	http.Handle("/statics/", http.StripPrefix("/statics/", fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on Port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
