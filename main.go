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

	userRoutes.HandleFunc("/home", handlers.Auth(handlers.UserHome))
	userRoutes.HandleFunc("/userValidate", handlers.ValidateUser)
	userRoutes.HandleFunc("/authenticate", handlers.SignUpPage)
	userRoutes.HandleFunc("/userRegister", handlers.UserRegister)
	userRoutes.HandleFunc("/logout", handlers.LogoutUser)

	adminRoutes := router.PathPrefix("/admin").Subrouter()

	adminRoutes.HandleFunc("/login", handlers.BackCheck(handlers.AdminLogin))
	adminRoutes.HandleFunc("/authenticate", handlers.BackCheck(handlers.ValidateAdmin))
	adminRoutes.HandleFunc("/dashboard", handlers.AdminAuth(handlers.Dashboard))
	adminRoutes.HandleFunc("/deleteUser/{Id}", handlers.AdminAuth(handlers.DeleteUser))
	adminRoutes.HandleFunc("/editUser/{Id}", handlers.AdminAuth(handlers.EditUser))
	adminRoutes.HandleFunc("/update/{Id}", handlers.AdminAuth(handlers.UpdateUser))
	adminRoutes.HandleFunc("/logout", handlers.AdminLogout)

	router.NotFoundHandler = http.HandlerFunc(handlers.NoRouteHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}
