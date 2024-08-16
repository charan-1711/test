package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"todo-list/database"
	"todo-list/models"
	"todo-list/utils"

	"github.com/golang-jwt/jwt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	user := new(models.User)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error while parsing json body", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	//validate request body
	if !utils.IsLength(user.Name, 5) {
		http.Error(w, "Username too short", http.StatusBadRequest)
		return
	}
	if !utils.IsLength(user.Password, 8) {
		http.Error(w, "Password must have atleast 8 characters", http.StatusBadRequest)
		return
	}

	storedUser := new(models.User)
	// get the registered user data
	err := database.DB.QueryRow("SELECT * FROM users WHERE username = $1", user.Name).
		Scan(&storedUser.UID, &storedUser.Name, &storedUser.Password)

	if err == sql.ErrNoRows { // check if user not found
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Password != storedUser.Password { // match the password
		http.Error(w, "Username or password is incorrect", http.StatusUnauthorized)
		return
	}

	// token creation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": storedUser.UID,
		"nbf":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create response body
	response := models.AuthResponse{Message: "User logged in successfully", Token: tokenString}
	jsonRes, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonRes)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	user := new(models.User)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error while parsing json body", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	//validate request body
	if !utils.IsLength(user.Name, 5) {
		http.Error(w, "Username too short", http.StatusBadRequest)
		return
	}
	if !utils.IsLength(user.Password, 8) {
		http.Error(w, "Password must have atleast 8 characters", http.StatusBadRequest)
		return
	}

	storedUser := new(models.User)
	database.DB.QueryRow("SELECT * FROM users WHERE username = $1", user.Name).
		Scan(&storedUser.UID, &storedUser.Name, &storedUser.Password)

	if storedUser.Name != "" { // check if user already exists
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// insert the new user details
	_, err := database.DB.Exec(
		"INSERT INTO users(username, password) values($1, $2)",
		user.Name,
		user.Password,
	)

	if err != nil {
		log.Print(err)
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: "User signup successfull"}
	jsonRes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}
