package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/caesar-john/task-management/repository"
)

func ToggleTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || taskID < 1 {
		http.Error(w, "Invalid task ID.", http.StatusBadRequest)
		return
	}

	status, err := repository.ToggleTaskStatus(r.Context(), taskID)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("Failed to toggle task status: %v", err)
		http.Error(w, "Could not update the task.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard?status="+status, http.StatusSeeOther)
}
