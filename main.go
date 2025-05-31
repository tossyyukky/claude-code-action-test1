package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskManager struct {
	db *sql.DB
}

func NewTaskManager() (*TaskManager, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		return nil, err
	}

	tm := &TaskManager{db: db}
	if err := tm.createTable(); err != nil {
		return nil, err
	}

	return tm, nil
}

func (tm *TaskManager) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		deadline DATETIME NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := tm.db.Exec(query)
	return err
}

func (tm *TaskManager) CreateTask(task Task) error {
	query := `
	INSERT INTO tasks (title, description, deadline, completed)
	VALUES (?, ?, ?, ?)`

	_, err := tm.db.Exec(query, task.Title, task.Description, task.Deadline, task.Completed)
	return err
}

func (tm *TaskManager) GetTasks() ([]Task, error) {
	query := `SELECT id, title, description, deadline, completed, created_at FROM tasks ORDER BY deadline ASC`
	rows, err := tm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var deadlineStr, createdAtStr string
		
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &deadlineStr, &task.Completed, &createdAtStr)
		if err != nil {
			return nil, err
		}

		task.Deadline, _ = time.Parse("2006-01-02 15:04:05", deadlineStr)
		task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tm *TaskManager) GetTask(id int) (Task, error) {
	query := `SELECT id, title, description, deadline, completed, created_at FROM tasks WHERE id = ?`
	row := tm.db.QueryRow(query, id)

	var task Task
	var deadlineStr, createdAtStr string
	
	err := row.Scan(&task.ID, &task.Title, &task.Description, &deadlineStr, &task.Completed, &createdAtStr)
	if err != nil {
		return task, err
	}

	task.Deadline, _ = time.Parse("2006-01-02 15:04:05", deadlineStr)
	task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	
	return task, nil
}

func (tm *TaskManager) UpdateTaskStatus(id int, completed bool) error {
	query := `UPDATE tasks SET completed = ? WHERE id = ?`
	_, err := tm.db.Exec(query, completed, id)
	return err
}

func (tm *TaskManager) Close() error {
	return tm.db.Close()
}

var taskManager *TaskManager
var templates *template.Template

func init() {
	var err error
	taskManager, err = NewTaskManager()
	if err != nil {
		log.Fatal("Failed to initialize task manager:", err)
	}

	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskManager.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Tasks []Task
	}{
		Tasks: tasks,
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")
		deadlineStr := r.FormValue("deadline")

		deadline, err := time.Parse("2006-01-02T15:04", deadlineStr)
		if err != nil {
			http.Error(w, "Invalid deadline format", http.StatusBadRequest)
			return
		}

		task := Task{
			Title:       title,
			Description: description,
			Deadline:    deadline,
			Completed:   false,
		}

		if err := taskManager.CreateTask(task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := templates.ExecuteTemplate(w, "create.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := taskManager.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func updateTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Completed bool `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := taskManager.UpdateTaskStatus(id, requestData.Completed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func main() {
	defer taskManager.Close()

	r := mux.NewRouter()
	
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/create", createTaskHandler).Methods("GET", "POST")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", getTaskHandler).Methods("GET")
	r.HandleFunc("/api/tasks/{id:[0-9]+}/status", updateTaskStatusHandler).Methods("POST")
	
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}