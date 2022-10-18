package handlers

import (
	"Login_Admin_Using_Postgres/database"
	"Login_Admin_Using_Postgres/models"
	"Login_Admin_Using_Postgres/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//............................Admin Handlers.............................................

//............................Admin Middleware.............................................

func AdminAuth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := utils.AdminStore.Get(r, "session")

		if session.Values["AdminAutheticate"] == false || session.Values["AdminAuthenticate"] == nil {
			val := models.Sample{ErrMsg: "You Must Login !", Title: "Admin Login"}
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
	val := models.Sample{Title: "Admin Login"}
	Tpl.ExecuteTemplate(w, "adminLogin.html", val)
}

func ValidateAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	err := r.ParseForm()
	if err != nil {
		log.Fatal("validate admin parse err : ", err)
	}
	userName := r.PostFormValue("adminName")
	userPass := r.PostFormValue("adminPassword")

	dbPass, emailValid := database.CheckAdminEmail(userName)
	if emailValid {
		if dbPass == userPass {
			fmt.Println("email valid")
			session, _ := utils.AdminStore.Get(r, "session")
			session.Values["AdminAuthenticate"] = true
			session.Values["AdminUserName"] = userName
			session.Save(r, w)
			http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
		} else {
			val := models.Sample{ErrMsg: "Invalid  Password", Title: "Admin Login"}
			Tpl.ExecuteTemplate(w, "adminLogin.html", val)
		}
	} else {
		val := models.Sample{ErrMsg: "Invalid Username", Title: "Admin Login"}
		Tpl.ExecuteTemplate(w, "adminLogin.html", val)
	}

}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.AdminStore.Get(r, "session")
	if session.Values["AdminAuthenticate"] == true {
		w.Header().Set("Cache-Control", "no-store")
		username := session.Values["AdminUserName"].(string)
		UserData := database.FetchUserData()
		UserData.Title = "DashBoard"
		UserData.AdminName = username
		Tpl.ExecuteTemplate(w, "adminDashboard.html", UserData)
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["Id"]
	res := database.DeleteUser(userId)
	if res {
		fmt.Print("worked")
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
	} else {
		fmt.Fprint(w, "not deleted!!")
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["Id"]
	fmt.Println("user id :", userId)
	email, err := database.GetUser(userId)
	if !err {
		log.Fatal(err)
	}
	session, _ := utils.AdminStore.Get(r, "session")
	userName := session.Values["AdminUserName"].(string)
	val := models.Sample{Title: "Edit User", Email: email, Id: userId, AdminName: userName}
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
	val := models.Sample{Title: "Admin Login"}
	Tpl.ExecuteTemplate(w, "adminLogin.html", val)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["Id"]
	r.ParseForm()
	newEmail := r.PostFormValue("updatedEmail")
	err := database.UpdateUserdata(userId, newEmail)
	if !err {
		log.Fatal("User update :", err)
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)

}

//............................No Route Handler.............................................

func NoRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
