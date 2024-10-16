package tasks

import (
	"database/sql"
	"fmt"
	"savaf/todo-list/utils"
	"time"
)

type Task struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type Tasks []Task

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB

// Initialize the database and create the "tasks" table if it doesn't exist
func InitializeDB(db *sql.DB) (*sql.DB, error) {
	DB = db
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		created_at TEXT NOT NULL,
		completed_at TEXT,
		is_deleted INTEGER DEFAULT 0
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		return nil, err
	}

	return DB, nil
}

// Show all Tasks
func ShowAll() (Tasks, error) {
	rows, err := DB.Query("SELECT id, title, created_at, completed_at FROM tasks WHERE is_deleted = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks Tasks
	for rows.Next() {
		var task Task
		var createdAt string
		var completedAt sql.NullString

		err := rows.Scan(&task.Id, &task.Title, &createdAt, &completedAt)
		if err != nil {
			return nil, err
		}

		// Parse time fields
		task.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			t, err := time.Parse(time.RFC3339, completedAt.String)
			if err != nil {
				return nil, err
			}
			task.CompletedAt = &t
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Show pending Tasks
func ShowPending() (Tasks, error) {
	rows, err := DB.Query("SELECT id, title, created_at, completed_at FROM tasks WHERE  completed_at IS NULL AND is_deleted = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks Tasks
	for rows.Next() {
		var task Task
		var createdAt string
		var completedAt sql.NullString

		err := rows.Scan(&task.Id, &task.Title, &createdAt, &completedAt)
		if err != nil {
			return nil, err
		}

		// Parse time fields
		task.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			t, err := time.Parse(time.RFC3339, completedAt.String)
			if err != nil {
				return nil, err
			}
			task.CompletedAt = &t
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func Add(title string) error {
	createdAt := time.Now().Format(time.RFC3339)
	query := `INSERT INTO tasks (title, created_at, completed_at) VALUES (?, ?, NULL)`
	_, err := DB.Exec(query, title, createdAt)
	if err != nil {
		return err
	}

	return nil
}

func Complete(taskId int64) error {
	completedAt := time.Now().Format(time.RFC3339)
	query := `UPDATE tasks SET completed_at = ? WHERE id = ?`
	_, err := DB.Exec(query, completedAt, taskId)
	if err != nil {
		return err
	}
	return nil
}

func GetById(taskId int64) (*Task, error) {
	var task Task
	var createdAt string
	var completedAt sql.NullString
	query := "SELECT id, title, created_at, completed_at FROM tasks WHERE id = ?"
	row := DB.QueryRow(query, taskId)
	err := row.Scan(&task.Id, &task.Title, &createdAt, &completedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no task found with ID %d", taskId)
		}
		return nil, err
	}

	// Parse time fields
	task.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		t, err := time.Parse(time.RFC3339, completedAt.String)
		if err != nil {
			return nil, err
		}
		task.CompletedAt = &t
	}

	return &task, nil
}

func Toggle(taskId int64) (*Task, error) {
	task, err := GetById(taskId)
	if err != nil {
		return nil, err
	}
	query := `UPDATE tasks SET completed_at = ? WHERE id = ?`
	var completedAt interface{}
	if task.CompletedAt == nil {
		completedAt = time.Now().Format(time.RFC3339)
	} else {
		completedAt = nil
	}
	_, err = DB.Exec(query, completedAt, taskId)
	if err != nil {
		return nil, err
	}

	task, err = GetById(taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func SoftDelete(taskId int64) (*Task, error) {
	task, err := GetById(taskId)
	if err != nil {
		return nil, err
	}
	query := `UPDATE tasks SET is_deleted = 1 WHERE id = ?`
	_, err = DB.Exec(query, taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func Delete(taskId int64) (*Task, error) {
	task, err := GetById(taskId)
	if err != nil {
		return nil, err
	}
	query := `DELETE FROM tasks WHERE id = ?`
	_, err = DB.Exec(query, taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (task Task) ToRow() []string {
	completionStatus := "Incompleted"
	if task.CompletedAt != nil {
		completionStatus = "Done"
	}
	return []string{
		fmt.Sprintf("%d", task.Id),
		task.Title,
		utils.TimeAgo(task.CreatedAt),
		completionStatus,
	}
}
