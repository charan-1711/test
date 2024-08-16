package main

import (
	"log"
	"net/http"
	"todo-list/controllers"
	"todo-list/database"
	"todo-list/middleware"

	_ "github.com/lib/pq"
)

func main() {
	database.ConnectToDB()

	http.HandleFunc("/auth/signup", controllers.Signup)
	http.HandleFunc("/auth/login", controllers.Login)
	http.Handle("/todo/create", middleware.IsAuth(http.HandlerFunc(controllers.CreateTodo)))
	http.Handle("/todo/update", middleware.IsAuth(http.HandlerFunc(controllers.UpdateTodo)))
	http.Handle("/todo/delete", middleware.IsAuth(http.HandlerFunc(controllers.DeleteTodo)))
	http.Handle("/todo/get", middleware.IsAuth(http.HandlerFunc(controllers.GetTodo)))
	http.Handle("/todo/getAll", middleware.IsAuth(http.HandlerFunc(controllers.GetAllTodo)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})

	err := http.ListenAndServe("localhost:3000", nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}
