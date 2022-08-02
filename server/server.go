package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wisuja/crud/service"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	valid, err := service.CheckLogin(r.PostForm.Get("username"), r.PostForm.Get("password"))

	if !valid {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())

		return
	}

	SetLoggedInCookie(w)
	fmt.Fprintln(w, "You've logged in!")
}

func FetchAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	users, err := service.FetchAllUsers()

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	for _, user := range users {
		fmt.Fprintln(w, user)
	}
}

func FetchUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Empty id")
		return
	}

	parsedId, _ := strconv.Atoi(id)
	user, err := service.FetchUser(parsedId)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, user)
}

func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	if !r.PostForm.Has("username") && !r.PostForm.Has("password") {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	user, err := service.CreateUser(r.PostForm.Get("username"), r.PostForm.Get("password"))

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Empty id")
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	if !r.PostForm.Has("username") && !r.PostForm.Has("password") {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	parsedId, _ := strconv.Atoi(id)
	user, err := service.UpdateUser(parsedId, r.PostForm.Get("username"), r.PostForm.Get("password"))

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Empty id")
		return
	}

	parsedId, _ := strconv.Atoi(id)
	_, err := service.DeleteUser(parsedId)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, "Successfully deleted id: ", id)
}

func CheckLoggedInCookie(w http.ResponseWriter, r *http.Request) bool {
	loggedIn, err := r.Cookie("X-WSJ-LOGGED-IN")

	if err != nil || loggedIn.Value != "true" {
		w.WriteHeader(400)
		fmt.Fprint(w, errors.New("You are not logged in!"))

		return false
	}

	return true
}

func SetLoggedInCookie(w http.ResponseWriter) {
	cookie := new(http.Cookie)
	cookie.Name = "X-WSJ-LOGGED-IN"
	cookie.Value = "true"
	cookie.Path = "/"

	http.SetCookie(w, cookie)
}

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login/", LoginHandler)
	mux.HandleFunc("/get-all-users/", FetchAllUsersHandler)
	mux.HandleFunc("/get-user/", FetchUserHandler)
	mux.HandleFunc("/insert-user/", InsertUserHandler)
	mux.HandleFunc("/update-user/", UpdateUserHandler)
	mux.HandleFunc("/delete-user/", DeleteUserHandler)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
