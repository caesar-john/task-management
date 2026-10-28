package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	User_id     int
	Title       string
	Description string
	status      string
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("dashboardHandler terpanggil")
	tpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/partials/header.html",
		"templates/partials/sidebar.html",
		"templates/pages/dashboard.html",
	)

	if err != nil {
		log.Println("ExecuteTemplate error: ", err)
		http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println("ExecuteTemplate error: ", err)
		http.Error(w, "Status Internal Server", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()

	//handler
	mux.HandleFunc("GET /dashboard", dashboardHandler)
	//buat css
	mux.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))

	var address = "localhost:8080"
	fmt.Printf("Server started at %s\n", address)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
