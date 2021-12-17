package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"
)

func DeleteMyself(w http.ResponseWriter, r *http.Request) {
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
