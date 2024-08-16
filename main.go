package main

import (
	"log"
	"net/http"
	"todo-list/connection"
	"todo-list/controllers"
	"todo-list/middleware"

	_ "github.com/lib/pq"
)

func main() {
	connection.ConnectToDB()

	http.HandleFunc("/auth/signup", controllers.Signup)
	http.HandleFunc("/auth/login", controllers.Login)
	http.Handle("/todo/create", middleware.IsAuth(http.HandlerFunc(controllers.CreateTodo))) // todos
	http.Handle("/todo/update", middleware.IsAuth(http.HandlerFunc(controllers.UpdateTodo))) // todos/:id method put
	http.Handle("/todo/delete", middleware.IsAuth(http.HandlerFunc(controllers.DeleteTodo))) // todos/:id method delete
	http.Handle("/todo/get", middleware.IsAuth(http.HandlerFunc(controllers.GetTodo))) // id  with get
	http.Handle("/todo/getAll", middleware.IsAuth(http.HandlerFunc(controllers.GetAllTodo))) // get
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})

	err := http.ListenAndServe("localhost:3000", nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}
