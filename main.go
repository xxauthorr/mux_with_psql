package main

import (
	"Login_Admin_Using_Postgres/database"
	"Login_Admin_Using_Postgres/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	var err error

	handlers.Tpl, err = handlers.Tpl.ParseGlob(("./static/*.html"))

	if err != nil {
		panic(err)
	}
	_, err1 := handlers.Tpl.New("partials").ParseGlob("./static/partials/*.html")

	if err1 != nil {
		log.Fatal(err1)
	}
}

func main() {
	database.ConnectDb()

	router := mux.NewRouter().StrictSlash(true)

	FileServer := http.FileServer(http.Dir("./static/assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", FileServer))

	router.HandleFunc("/", handlers.IndexPage)
	userRoutes := router.PathPrefix("/user").Subrouter()

	userRoutes.HandleFunc("/authenticate", handlers.SignUpPage)
	userRoutes.HandleFunc("/home", handlers.Auth(handlers.UserHome))
	userRoutes.HandleFunc("/userRegister", handlers.UserRegister)
	userRoutes.HandleFunc("/userValidate", handlers.ValidateUser)
	userRoutes.HandleFunc("/logout", handlers.LogoutUser)

	log.Fatal(http.ListenAndServe(":3000", router))
}
