package handlers

import (
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte("this-is-secret"))

func HomePage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println(err)                              // Ugly debug output
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	//we send a message if the user is connected so that some buttons will not display
	var message structs.Comment
	if !session.IsNew {
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", nil)
	}
}

func AboutUs(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println(err)                              // Ugly debug output
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	//we send a message if the user is connected so that some buttons will not display
	var message structs.Comment
	if !session.IsNew {
		message.ID = "yes"
		temp.ExecuteTemplate(w, "about.html", message)
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "about.html", nil)
	}
}

func PersonalPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	//we send a message if the user is connected so that some buttons will not display
	var message structs.Comment
	if !session.IsNew {
		message.ID = "yes"
		if session.Values["Role"].(string) == "AGENCY" {
			message.IsAgency = "yes"
			message.Username = session.Values["Username"].(string)
			message.IsTheSame = "yes"
			message.ID = "yes"
			temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
			return
		}
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		if session.Values["Role"].(string) == "CLIENT" {
			//take the title
			sqlQueryA := "SELECT money_balance from clients where id=?"
			rowA := database.Db.QueryRow(sqlQueryA, session.Values["Id"].(string))
			var money_balance string
			//if don't exist => no money_balance
			if rowA.Scan(&money_balance) != nil {
				//we send a message if the user is connected so that some buttons will not display
				var message structs.Comment
				if !session.IsNew {
					message.ID = "yes"
					message.ErrMessage = "Ceva nu a mers cum trebuie!"
					temp.ExecuteTemplate(w, "index.html", message)
				} else {
					session.Options.MaxAge = -1
					session.Save(r, w)
					message.ErrMessage = "Ceva nu a mers cum trebuie!"
					temp.ExecuteTemplate(w, "index.html", message)
				}
				return
			}
			message.MoneyBalance = money_balance
		}
		temp.ExecuteTemplate(w, "personalPage.html", message)
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", nil)
	}
}

func Login(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	var message structs.Comment
	if !session.IsNew {
		message.ID = "yes"
		if session.Values["Role"].(string) == "AGENCY" {
			message.IsAgency = "yes"
			message.Username = session.Values["Username"].(string)
			message.IsTheSame = "yes"
			message.ID = "yes"
			temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
			return
		}
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		message.Username = "Esti deja conectat"
		temp.ExecuteTemplate(w, "personalPage.html", message)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "login.html", nil)
}

func LoginLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	// session, _ := Store.Get(r, "session")
	// var message structs.Comment
	// if !session.IsNew {
	// 	message.ID = "yes"
	// 	if session.Values["Role"].(string) == "AGENCY" {
	// 		message.IsAgency = "yes"
	// 		message.Username = session.Values["Username"].(string)
	// 		message.IsTheSame = "yes"
	// 		temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
	// 		return
	// 	}
	// 	if session.Values["Role"].(string) == "ADMIN" {
	// 		message.IsAdmin = "yes"
	// 	}
	// 	message.Username = "Esti deja conectat"
	// 	temp.ExecuteTemplate(w, "personalPage.html", message)
	// 	return
	// }
	// session.Options.MaxAge = -1
	// session.Save(r, w)
	err1 := r.ParseForm()

	if err1 != nil {
		return
	}
	// select the username and password
	username := r.FormValue("username")
	password := r.FormValue("password")

	//validate username and password
	err1 = checkUsername(username)
	if err1 != nil {
		var message structs.Comment
		message.Message = "Numele de utilizator sau parola nu sunt corecte!"
		temp.ExecuteTemplate(w, "login.html", message)
		return
	}

	//same for password
	err1 = checkPassword(password)
	if err1 != nil {
		var message structs.Comment
		message.Message = "Numele de utilizator sau parola nu sunt corecte!"
		temp.ExecuteTemplate(w, "login.html", message)
		return
	}

	sqlQuery := "SELECT id,hash,role from clients where username=?"
	row := database.Db.QueryRow(sqlQuery, username)
	var id, hash, role string
	//select in another query the particularities from agencies
	sqlQueryA := "SELECT id,hash,role from agencies where username=?"
	rowA := database.Db.QueryRow(sqlQueryA, username)
	var idA, hashA, roleA string
	//if both don't exist => no user
	if row.Scan(&id, &hash, &role) != nil && rowA.Scan(&idA, &hashA, &roleA) != nil {
		var message structs.Comment
		message.Message = "Numele de utilizator sau parola nu sunt corecte!"
		temp.ExecuteTemplate(w, "login.html", message)
		return
	}

	//we don't let the user log in if he didn't verify the account
	sqlQuery2 := "SELECT is_active from clients where username=?"
	row2 := database.Db.QueryRow(sqlQuery2, username)
	var is_active string
	//we do the same thing for agencies
	sqlQuery2A := "SELECT is_active from agencies where username=?"
	row2A := database.Db.QueryRow(sqlQuery2A, username)
	var is_activeA string
	if row2.Scan(&is_active) != nil && row2A.Scan(&is_activeA) != nil {
		var message structs.Comment
		message.Message = "Numele de utilizator sau parola nu sunt corecte!"
		temp.ExecuteTemplate(w, "login.html", message)
		return
	}

	if is_active == "0" {
		var message structs.Comment
		message.ErrMessage = "Va rugam sa verificati emailul inainte de a va conecta"
		temp.ExecuteTemplate(w, "emailVerification.html", message)
		return
	}
	var err error
	// verify the password
	if roleA == "AGENCY" {
		err = bcrypt.CompareHashAndPassword([]byte(hashA), []byte(password))
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	}
	if err == nil {
		session, _ := Store.Get(r, "session")
		if roleA == "AGENCY" {
			session.Values["Id"] = idA
			session.Values["Role"] = roleA
			session.Values["Username"] = username
		} else {
			session.Values["Id"] = id
			session.Values["Role"] = role
			session.Values["Username"] = username
		}
		fmt.Println(roleA)
		session.Save(r, w)
		var message structs.Comment
		if session.Values["Role"].(string) == "AGENCY" {
			message.IsAgency = "yes"
			message.Username = session.Values["Username"].(string)
			message.ErrMessage = "Esti conectat ca" + username
			message.IsTheSame = "yes"
			message.ID = "yes"
			temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
			return
		}
		message.Username = "Esti conectat ca " + username
		message.ID = "yes"
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		temp.ExecuteTemplate(w, "personalPage.html", message)
		return
	}
	var message structs.Comment
	message.Message = "Numele de utilizator sau parola nu sunt corecte!"
	temp.ExecuteTemplate(w, "login.html", message)

}

func Logout(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	//close the session
	temp, _ := template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	delete(session.Values, "Id")
	session.Options.MaxAge = -1
	session.Save(r, w)
	var message structs.Comment
	message.Message = "Ai fost deconectat!"
	temp.ExecuteTemplate(w, "login.html", message)
}
