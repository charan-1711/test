package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-list/database"
	"todo-list/models"
	"todo-list/utils"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	todo := new(models.Todo)
	err := json.NewDecoder(r.Body).Decode(&todo)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error while parsing json body", http.StatusUnprocessableEntity)
		return
	}

	//validate request body
	if !utils.IsLength(todo.Title, 5) {
		http.Error(w, "Todo title too short", http.StatusBadRequest)
		return
	}
	if !utils.Contains([]string{"Not Started", "In Progress", "Completed"}, todo.Status) {
		http.Error(w, "Invalid progress status", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value("userId").(int)

	_, err = database.DB.Exec(
		"INSERT INTO todos(title, current_status, u_id) VALUES($1, $2, $3)",
		todo.Title,
		todo.Status,
		userId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(models.Response{Message: "Todo created successfully"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	todo := new(models.Todo)
	err := json.NewDecoder(r.Body).Decode(&todo)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if todo.TID == 0 {
		http.Error(w, "Invalid Todo Id", http.StatusBadRequest)
		return
	}

	var storedUserId int
	err = database.DB.QueryRow("SELECT u_id FROM todos WHERE t_id = $1", todo.TID).
		Scan(&storedUserId)
	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get requested user's userId
	userId := r.Context().Value("userId").(int)

	// check to match the userId
	if userId != storedUserId {
		http.Error(w, "Your not authorized to delete this todo", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("DELETE FROM todos WHERE t_id = $1", todo.TID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(models.Response{Message: "Succesfully deleted todo"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	todo := new(models.Todo)
	err := json.NewDecoder(r.Body).Decode(&todo)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error while parsing json body", http.StatusUnprocessableEntity)
		return
	}

	//validate request body
	if !utils.IsLength(todo.Title, 5) {
		http.Error(w, "Todo title too short", http.StatusBadRequest)
		return
	}
	if !utils.Contains([]string{"Not Started", "In Progress", "Completed"}, todo.Status) {
		http.Error(w, "Invalid progress status", http.StatusBadRequest)
		return
	}

	var storedUserId int
	err = database.DB.QueryRow("SELECT u_id FROM todos WHERE t_id = $1", todo.TID).
		Scan(&storedUserId)

	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := r.Context().Value("userId").(int)

	//	check to match the userId
	if userId != storedUserId {
		http.Error(w, "Your not authorized to update this todo", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec(
		"UPDATE todos SET title = $1, current_status = $2 WHERE t_id = $3",
		todo.Title,
		todo.Status,
		todo.TID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(models.Response{Message: "Todo updated successfully"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	t_id, _ := strconv.Atoi(r.URL.Query().Get("t_id"))
	todo := new(models.Todo)
	var storedUserId int
	err := database.DB.QueryRow("SELECT t_id, title, current_status, u_id FROM todos WHERE t_id = $1", t_id).
		Scan(&todo.TID, &todo.Title, &todo.Status, &storedUserId)

	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := r.Context().Value("userId").(int)
	if userId != storedUserId {
		http.Error(w, "Your not authorized to view this todo", http.StatusForbidden)
		return
	}

	jsonRes, err := json.Marshal(todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}

func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userId").(int)

	rows, err := database.DB.Query(
		"SELECT t_id, title, current_status FROM todos WHERE u_id = $1",
		userId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result = []*models.Todo{}

	for rows.Next() {
		todo := new(models.Todo)
		rows.Scan(&todo.TID, &todo.Title, &todo.Status)

		result = append(result, todo)
	}

	jsonRes, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}
