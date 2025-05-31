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

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		deadline DATETIME,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	
	rows, err := db.Query("SELECT id, title, description, deadline, completed, created_at FROM tasks ORDER BY deadline ASC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var deadlineStr string
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &deadlineStr, &task.Completed, &task.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		task.Deadline, _ = time.Parse("2006-01-02 15:04:05", deadlineStr)
		tasks = append(tasks, task)
	}

	tmpl.Execute(w, tasks)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/create.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")
		deadlineStr := r.FormValue("deadline")
		
		deadline, err := time.Parse("2006-01-02T15:04", deadlineStr)
		if err != nil {
			http.Error(w, "Invalid deadline format", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO tasks (title, description, deadline) VALUES (?, ?, ?)",
			title, description, deadline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func toggleTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET completed = NOT completed WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task Task
	var deadlineStr string
	err = db.QueryRow("SELECT id, title, description, deadline, completed, created_at FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.Title, &task.Description, &deadlineStr, &task.Completed, &task.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.Deadline, _ = time.Parse("2006-01-02 15:04:05", deadlineStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/create", createTaskHandler)
	r.HandleFunc("/task/{id}/toggle", toggleTaskHandler).Methods("POST")
	r.HandleFunc("/task/{id}", getTaskHandler).Methods("GET")
	
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Println("サーバーを起動しています: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}