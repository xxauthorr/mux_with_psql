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

//........................User MiddleWares........................................

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

//........................User Handlers.......................................

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
	session, _ := utils.UserStore.Get(r, "session")
	if session.Values["authenticated"] == true {
		http.Redirect(w, r, "/user/home", http.StatusFound)
	}
	val := models.Credentials{Header: "Sign up"}
	Tpl.ExecuteTemplate(w, "login.html", val)

}

// .........................Login Check.......................................
func ValidateUser(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Fatal("parse Err : ", err)
	}

	email := r.PostFormValue("userEmail")
	pass := r.PostFormValue("userPass")
	fmt.Println(email)
	hashPass, emailValid := database.CheckUserEmail(email)
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

//.........................User Signup.......................................

func UserRegister(w http.ResponseWriter, r *http.Request) {
	for err := r.ParseForm(); err != nil; {
		log.Fatal(err)
	}

	email := r.PostFormValue("userEmail")
	pass := r.PostFormValue("userPassConfirm")
	sessions, _ := utils.UserStore.Get(r, "session")

	_, emailFound := database.CheckUserEmail(email)
	if emailFound {
		w.Header().Set("Cache-Control", "no-store")
		val := models.Credentials{Header: "Login", Errmsg: "Email already exist"}
		Tpl.ExecuteTemplate(w, "login.html", val)
	}

	sessions.Values["authenticated"] = true
	sessions.Values["EmailId"] = email
	sessions.Save(r, w)
	database.RegisterUser(email, pass)
	http.Redirect(w, r, "/user/home", http.StatusFound)
}

//..........................Login to Home..............................................

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

// ..............................Logout User.............................................
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

//............................Admin Handlers.............................................

func AdminAuth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := utils.AdminStore.Get(r, "session")

		if session.Values["AdminAutheticate"] == false || session.Values["AdminAuthenticate"] == nil {
			val := models.ClientUser{ErrMsg: "You Must Login !", Title: "Admin Login"}
			Tpl.ExecuteTemplate(w, "adminLogin.html", val)
			return
		}
		HandlerFunc.ServeHTTP(w, r)
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.AdminStore.Get(r, "session")
	if session.Values["AdminAutheticate"] == true {
		http.Redirect(w, r, "/admin/dasboard", http.StatusFound)
	}
	fmt.Println("AdminLogin")
	val := models.ClientUser{Title: "Admin Login"}
	Tpl.ExecuteTemplate(w, "adminLogin.html", val)
}

func ValidateAdmin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("validate admin parse err : ", err)
	}
	userName := r.PostFormValue("adminName")
	// password := r.PostFormValue("adminPassword")

	_, emailValid := database.CheckAdminEmail(userName)
	if emailValid {
		// passValid := utils.CheckPasswordMatch(password, hashPass)
		// if passValid {
		fmt.Println("email valid")
		w.Header().Set("Cache-Control", "no-store")
		session, _ := utils.AdminStore.Get(r, "session")
		session.Values["AdminAuthenticate"] = true
		session.Values["AdminUserName"] = userName
		session.Save(r, w)
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
		// } else {
		// 	w.Header().Set("Cache-Control", "no-store")
		// 	val := models.Credentials{Errmsg: "Wrong PassWord"}
		// 	Tpl.ExecuteTemplate(w, "adminLogin.html", val)
		// }
	} else {
		w.Header().Set("Cache-Control", "no-store")
		val := models.ClientUser{ErrMsg: "Invalid Username", Title: "Admin Login"}
		Tpl.ExecuteTemplate(w, "adminLogin.html", val)
	}

}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.AdminStore.Get(r, "session")
	if session.Values["AdminAuthenticate"] == true {
		w.Header().Set("Cache-Control", "no-store")
		// username := session.Values["AdminUserName"].(string)
		UserData := database.FetchUserData()
		// Data := models.ClientUser{Title: "Dasboard", Username: username}
		// UserData = append(UserData, Data)
		// for i := range UserData {
		// 	models.ClientUser{Id: UserData[i].Id}
		// }
		UserData.Title = "Dash"
		Tpl.ExecuteTemplate(w, "adminDashboard.html", UserData)
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostFormValue("userId")
	database.DeleteUser(id)
}

func EditUser(w http.ResponseWriter, r *http.Request) {

	val := models.ClientUser{Title: "Edit User"}
	Tpl.ExecuteTemplate(w, "editPanel.html", val)

}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	session, _ := utils.AdminStore.Get(r, "session")
	if session.Values["AdminAuthenticate"] == false {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	session.Values["AdminUserName"] = nil
	session.Values["AdminAuthenticate"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)
	val := models.ClientUser{Title: "Admin Login"}
	Tpl.ExecuteTemplate(w, "adminLogin.html", val)

}
