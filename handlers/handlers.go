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

// Authenticate middleware
func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := utils.UserStore.Get(r, "session")

		if session.Values["autheticated"] == false || session.Values["authenticated"] == nil {
			val := models.Credentials{Header: "Login", Errmsg: "You Must Login !"}
			Tpl.ExecuteTemplate(w, "login.html", val)
			return
		}
		HandlerFunc.ServeHTTP(w, r)
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.UserStore.Get(r, "session")
	if session.Values["authenticated"] == true {
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, "/user/home", http.StatusFound)
	}
	val := models.Credentials{Header: "Home"}
	Tpl.ExecuteTemplate(w, "guest.html", val)
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
			w.Header().Set("Cache-Control", "no-store")
			session, _ := utils.UserStore.Get(r, "session")
			session.Values["authenticated"] = true
			session.Values["EmailId"] = email
			session.Save(r, w)
			http.Redirect(w, r, "/user/home", http.StatusFound)
		} else {
			fmt.Println("invalid Pass")
			w.Header().Set("Cache-Control", "no-store")
			val := models.Credentials{Header: "Login", Errmsg: "Incorrect Password"}
			Tpl.ExecuteTemplate(w, "login.html", val)
		}
	} else {
		fmt.Println("invalid email")
		w.Header().Set("Cache-Control", "no-store")
		val := models.Credentials{Errmsg: "Invalid Email", Header: "Login"}
		Tpl.ExecuteTemplate(w, "login.html", val)
	}

}

func UserHome(w http.ResponseWriter, r *http.Request) {
	session, err := utils.UserStore.Get(r, "session")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Cache-Control", "no-store")
	fmt.Println(session.Values["authenticated"], session.Values["EmailId"])

	if session.Values["authenticated"] == true {
		email := session.Values["EmailId"].(string)
		val := models.Credentials{Header: "Home Page", Email: email}
		Tpl.ExecuteTemplate(w, "loggedIn.html", val)
	} else {
		val := models.Credentials{Errmsg: "You Must Login", Header: "Login"}
		Tpl.ExecuteTemplate(w, "login.html", val)
	}

}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	for err := r.ParseForm(); err != nil; {
		log.Fatal(err)
	}

	email := r.PostFormValue("userEmail")
	pass := r.PostFormValue("userPassConfirm")
	sessions, _ := utils.UserStore.Get(r, "session")

	_, emailValid := database.CheckEmail(email)
	if !emailValid {
		val := models.Credentials{Header: "Login", Errmsg: "Email already exist"}
		Tpl.ExecuteTemplate(w, "login.html", val)
	}

	sessions.Values["authenticated"] = true
	sessions.Values["EmailId"] = email
	sessions.Save(r, w)
	database.RegisterUser(email, pass)
	http.Redirect(w, r, "/user/home", http.StatusFound)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Clear the cache
	w.Header().Set("Cache-Control", "no-store")
	session, _ := utils.UserStore.Get(r, "session")
	session.Values["authenticated"] = nil
	session.Values["EmailId"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
