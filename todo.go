package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	ID        int
	Task      string
	Completed bool
}

var todos []Todo
var nextID int

func main() {
	// Wrap our handlers with the logging middleware
	http.HandleFunc("/", logRequest(listTodos))
	http.HandleFunc("/add", logRequest(addTodo))
	http.HandleFunc("/toggle", logRequest(toggleTodo))
	http.HandleFunc("/delete", logRequest(deleteTodo))

	fmt.Println("Server is running on http://localhost:8069")
	log.Fatal(http.ListenAndServe(":8069", nil))
}

// Middleware to log requests
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func listTodos(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Todo List</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        ul { list-style-type: none; padding: 0; }
        li { margin-bottom: 10px; }
        .completed { text-decoration: line-through; }
    </style>
</head>
<body>
    <h1>Todo List</h1>
    <form action="/add" method="post">
        <input type="text" name="task" required>
        <input type="submit" value="Add Todo">
    </form>
    <ul>
        {{range .}}
            <li class="{{if .Completed}}completed{{end}}">
                {{.Task}}
                <form style="display: inline;" action="/toggle" method="post">
                    <input type="hidden" name="id" value="{{.ID}}">
                    <input type="submit" value="Toggle">
                </form>
                <form style="display: inline;" action="/delete" method="post">
                    <input type="hidden" name="id" value="{{.ID}}">
                    <input type="submit" value="Delete">
                </form>
            </li>
        {{end}}
    </ul>
</body>
</html>
`
	t, _ := template.New("todos").Parse(tmpl)
	t.Execute(w, todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	task := r.FormValue("task")
	if task != "" {
		todos = append(todos, Todo{ID: nextID, Task: task, Completed: false})
		nextID++
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func toggleTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = !todos[i].Completed
			break
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}