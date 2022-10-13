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

// func FetchUserData() {
// 	ClientData := make([]*models.ClientUser, 0)
// 	rows, err := database.Db.Query(`SELECT id,email FROM clientuser`)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		user := new(models.ClientUser)
// 		if err := rows.Scan(&user.Id, &user.Email); err != nil {
// 			log.Fatal(err.Error())
// 		}
// 		ClientData = append(ClientData, user)
// 	}
// 	// fmt.Println(ClientData, "hai")
// 	for i := range ClientData {
// 		fmt.Println(*ClientData[i])
// 	}

// }

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

	adminRoutes.HandleFunc("/login", handlers.AdminLogin)
	adminRoutes.HandleFunc("/authenticate", handlers.ValidateAdmin)
	adminRoutes.HandleFunc("/dashboard", handlers.AdminAuth(handlers.Dashboard))
	adminRoutes.HandleFunc("/deleteUser/{Id}", handlers.AdminAuth(handlers.DeleteUser))
	adminRoutes.HandleFunc("/editUser/{Id}", handlers.AdminAuth(handlers.EditUser))
	adminRoutes.HandleFunc("/update/{Id}", handlers.AdminAuth(handlers.UpdateUser))
	adminRoutes.HandleFunc("/logout", handlers.AdminLogout)

	log.Fatal(http.ListenAndServe(":3000", router))
}
