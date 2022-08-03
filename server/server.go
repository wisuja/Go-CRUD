package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wisuja/crud/database"
	"github.com/wisuja/crud/entity"
	"github.com/wisuja/crud/repository"
)

type GlobalStruct struct {
	DB             *sql.DB
	UserRepository repository.UserRepository
}

func (gs *GlobalStruct) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	ctx := context.Background()

	_, err := gs.UserRepository.FindByUser(ctx, entity.User{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	})

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
	}

	SetLoggedInCookie(w)
	fmt.Fprintln(w, "You've logged in!")
}

func (gs *GlobalStruct) fetchAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	ctx := context.Background()
	users, err := gs.UserRepository.FindAll(ctx)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	for _, user := range users {
		fmt.Fprintln(w, user)
	}
}

func (gs *GlobalStruct) fetchUserHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx := context.Background()

	parsedId, _ := strconv.Atoi(id)
	user, err := gs.UserRepository.FindById(ctx, parsedId)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, user)
}

func (gs *GlobalStruct) insertUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "Invalid request method")
		return
	}

	loggedIn := CheckLoggedInCookie(w, r)
	if !loggedIn {
		return
	}

	if r.PostFormValue("username") == "" && r.PostFormValue("password") == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	ctx := context.Background()

	user, err := gs.UserRepository.Insert(ctx, entity.User{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	})

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, user)
}

func (gs *GlobalStruct) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if r.PostFormValue("username") == "" && r.PostFormValue("password") == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Input is invalid")

		return
	}

	ctx := context.Background()

	parsedId, _ := strconv.Atoi(id)
	user, err := gs.UserRepository.Update(ctx, parsedId, entity.User{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	})

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, user)
}

func (gs *GlobalStruct) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx := context.Background()

	parsedId, _ := strconv.Atoi(id)
	_, err := gs.UserRepository.Delete(ctx, parsedId)

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
	credentials, err := database.GetDefaultDatabaseConfig(".env")
	if err != nil {
		panic(err)
	}

	db, err := database.GetConnection(credentials)
	if err != nil {
		panic(err)
	}

	globalStruct := GlobalStruct{
		DB:             db,
		UserRepository: repository.NewUserRepository(db),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login/", globalStruct.loginHandler)
	mux.HandleFunc("/get-all-users/", globalStruct.fetchAllUsersHandler)
	mux.HandleFunc("/get-user/", globalStruct.fetchUserHandler)
	mux.HandleFunc("/insert-user/", globalStruct.insertUserHandler)
	mux.HandleFunc("/update-user/", globalStruct.updateUserHandler)
	mux.HandleFunc("/delete-user/", globalStruct.deleteUserHandler)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
