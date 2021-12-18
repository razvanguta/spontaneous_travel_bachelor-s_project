package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func DeleteMyself(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", "Ceva nu a functionat cum trebuie")
		return
	}

	var trans *sql.Tx
	trans, err := database.Db.Begin()

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "A aparut o eroare, te rog mai incearca!")
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	var deleteCustomer *sql.Stmt
	var deleteAgency *sql.Stmt

	if session.Values["Role"].(string) == "AGENCY" {
		deleteAgency, err = trans.Prepare("DELETE FROM agencies where id=?")
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut sterge!"
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer deleteAgency.Close()
		_, err = deleteAgency.Exec(session.Values["Id"].(string))
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut sterge!"
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer deleteAgency.Close()

	} else {
		deleteCustomer, err = trans.Prepare("DELETE FROM clients where id=?")
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut sterge!"
			temp.ExecuteTemplate(w, "index.html", message)
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer deleteCustomer.Close()
		_, err = deleteCustomer.Exec(session.Values["Id"].(string))
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut sterge!"
			temp.ExecuteTemplate(w, "index.html", message)
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer deleteCustomer.Close()
	}
	delete(session.Values, "Id")
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "login.html", nil)
	trans.Commit()

}

func EditDescriptionAgency(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", "Ceva nu a functionat cum trebuie")
		return
	}
	// we check if we are connected with an agency account
	if session.Values["Role"] != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu se poate accesa!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	temp.ExecuteTemplate(w, "editDescriptionAgency.html", nil)
}

func EditDescriptionAgencyPut(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	// we check if we are connected with an agency account
	if session.Values["Role"] != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu se poate accesa!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", "Ceva nu a functionat cum trebuie")
		return
	}
	//start the transaction
	var trans *sql.Tx
	trans, err := database.Db.Begin()

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "A aparut o eroare, te rog mai incearca!")
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	//we start the porccess of description update
	var updateDescription *sql.Stmt
	updateDescription, err = trans.Prepare("UPDATE agencies SET description=? where username=?")
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut edita!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer updateDescription.Close()
	err1 := r.ParseForm()

	if err1 != nil {
		return
	}
	fmt.Println(r)
	//execute with the new description and the username
	_, err = updateDescription.Exec(r.FormValue("description"), session.Values["Username"].(string))
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut edita!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer updateDescription.Close()
	//return to the agency page
	var message structs.Comment
	message.IsAgency = "yes"
	message.Username = session.Values["Username"].(string)
	message.IsTheSame = "yes"
	temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
	trans.Commit()
}
