package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Пишем 1 обработчик ( для получения всех задач)
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		log.Fatalf("error marshal tasks %v", err)
		http.Error(w, "error marshal tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)

	if err != nil {
		log.Fatalf("error writing response: %v", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}

}

// Пишем 2 обработчик (для отправки задачи на сервер)
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Fatalf("error reading body: %v", err)
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Fatalf("error unmarshalling body: %v", err)
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	for _, t := range tasks {
		if t.ID == task.ID {
			log.Fatalf("task %v already exists", t.ID)
			http.Error(w, "task already exists", http.StatusBadRequest)
			return
		}
	}

	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Пишем 3 обработчки для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		log.Fatalf("task %s not found", id)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		log.Fatalf("error marshal task %v", err)
		http.Error(w, "error marshal task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Fatalf("error writing response: %v", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}
}

// Пишем 4 обработчик удаления задачи по ID
func delTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := tasks[id]; !ok {
		http.Error(w, "Задача с таким ID не найдена", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Fatalf("error writing response: %v", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()
	// Регистрируем первый обработчик
	r.Get("/tasks", getTasks)
	// Регистрируем второй обработчик
	r.Post("/tasks", postTask)
	// Регистрируем третий обработчик
	r.Get("/tasks/{id}", getTask)
	// Регистрируем четвертый обработчик
	r.Delete("/tasks/{id}", delTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
