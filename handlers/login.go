package handlers

import (
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte("this-is-secret"))

func HomePage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	//we send a message if the user is connected so that some buttons will not display
	var message structs.Comment
	if !session.IsNew {
		message.ClientID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", nil)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		temp.ExecuteTemplate(w, "personalPage.html", "Esti deja conectat")
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "login.html", nil)
}

func LoginLogic(w http.ResponseWriter, r *http.Request) {
	temp, _ = template.ParseGlob("templates/*.html")
	err1 := r.ParseForm()

	if err1 != nil {
		return
	}
	// select the username and password
	username := r.FormValue("username")
	password := r.FormValue("password")

	sqlQuery := "SELECT id,hash from clients where username=?"
	row := database.Db.QueryRow(sqlQuery, username)
	var id, hash string
	if row.Scan(&id, &hash) != nil {
		temp.ExecuteTemplate(w, "login.html", "Numele de utilizator este incorect introdus")
		return
	}

	//we don't let the user log in if he didn't verify the account
	sqlQuery2 := "SELECT is_active from clients where username=?"
	row2 := database.Db.QueryRow(sqlQuery2, username)
	var is_active string
	if row2.Scan(&is_active) != nil {
		temp.ExecuteTemplate(w, "login.html", "Numele de utilizator este incorect introdus")
		return
	}

	if is_active == "0" {
		temp.ExecuteTemplate(w, "emailVerification.html", "Va rugam sa verificati emailul inainte de a va conecta")
		return
	}

	// verify the password

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		session, _ := Store.Get(r, "session")
		session.Values["clientId"] = id
		session.Save(r, w)
		temp.ExecuteTemplate(w, "personalPage.html", "Esti conectat ca "+username)
		return
	}

	temp.ExecuteTemplate(w, "login.html", "Parola este incorect introdusa")

}

func Logout(w http.ResponseWriter, r *http.Request) {
	//close the session
	temp, _ := template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	delete(session.Values, "id")
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "login.html", "Ai fost deconectat")
}
