package entities

import (
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Deadline    time.Time
	CreatedAt   time.Time
}

type DashboardViewData struct {
	Tasks          []Task
	TotalTasks     int
	PendingTasks   int
	CompletedTasks int
	Filter         string
	DateLabel      string
	Greeting       string
}
