package handlers

import (
	"Login_Admin_Using_Postgres/database"
	"Login_Admin_Using_Postgres/models"
	"Login_Admin_Using_Postgres/utils"
	"fmt"
	"html/template"

	"log"
	"net/http"
)

var Tpl *template.Template

func IndexPage(w http.ResponseWriter, r *http.Request) {
	val := models.Credentials{Header: "Home"}
	Tpl.ExecuteTemplate(w, "home.html", val)
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	val := models.Credentials{Header: "Sign up"}
	Tpl.ExecuteTemplate(w, "login.html", val)

}

func ValidateUser(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Fatal("parse Err : ", err)
	}

	email := r.PostFormValue("userEmail")
	pass := r.PostFormValue("userPass")
	fmt.Println(email)
	hashPass, emailValid := database.CheckEmail(email)
	if emailValid {
		fmt.Println("emailValid")
		passValid := utils.CheckPasswordMatch(pass, hashPass)
		if passValid {
			fmt.Println("passValid")
			session, _ := utils.UserStore.Get(r, "session")
			session.Values["authrnticated"] = true
			session.Values["EmailId"] = email
			val := models.Credentials{Header: "Home"}
			Tpl.ExecuteTemplate(w, "home.html", val)
		} else {
			fmt.Println("invalid Pass")
			val := models.Credentials{Header: "Login", Errmsg: "Incorrect Password"}
			Tpl.ExecuteTemplate(w, "login.html", val)
		}
	} else {
		fmt.Println("invalid email")
		val := models.Credentials{Errmsg: "Invalid Email", Header: "Login"}
		Tpl.ExecuteTemplate(w, "login.html", val)
	}

}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	for err := r.ParseForm(); err != nil; {
		log.Fatal(err)
	}

	email := r.PostFormValue("userEmail")
	pass := r.PostFormValue("userPassConfirm")

	database.RegisterUser(email, pass)
	val := models.Credentials{Email: email, Header: "Home"}
	Tpl.ExecuteTemplate(w, "home.html", val)
}
