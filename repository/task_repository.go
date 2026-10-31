package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/caesar-john/task-management/config"
	"github.com/caesar-john/task-management/entities"
)

func GetAllTask(ctx context.Context) ([]entities.Task, error) {
	row, err := config.DB.QueryContext(ctx,
		"SELECT id, title, description, status, deadline, created_at FROM tasks")

	if err != nil {
		log.Println("Failed to get all task: ", err)
	}

	var hasil []entities.Task

	for row.Next() {
		var t entities.Task

		if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Deadline, &t.CreatedAt); err != nil {
			return nil, err
		}

		hasil = append(hasil, t)

	}

	if err := row.Err(); err != nil {
		return nil, err
	}
	return hasil, err
}

func CreateTask(ctx context.Context, task entities.Task) error {
	_, err := config.DB.ExecContext(ctx,
		"INSERT INTO tasks (title, description, status, deadline, created_at) VALUES (?, ?, ?, ?, ?)",
		task.Title,
		task.Description,
		task.Status,
		task.Deadline,
		task.CreatedAt,
	)
	return err
}

func ToggleTaskStatus(ctx context.Context, taskID int) (string, error) {
	tx, err := config.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var currentStatus string
	if err := tx.QueryRowContext(ctx,
		"SELECT status FROM tasks WHERE id = ? FOR UPDATE", taskID,
	).Scan(&currentStatus); err != nil {
		return "", err
	}

	newStatus := "Completed"
	if strings.EqualFold(currentStatus, "Completed") {
		newStatus = "Pending"
	}

	result, err := tx.ExecContext(ctx,
		"UPDATE tasks SET status = ? WHERE id = ?", newStatus, taskID,
	)
	if err != nil {
		return "", err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return newStatus, nil
}
