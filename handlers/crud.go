package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"myapp/database"
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

	var insertCustomer *sql.Stmt
	insertCustomer, err = trans.Prepare("DELETE FROM clients where id=?")

	if err != nil {
		temp.ExecuteTemplate(w, "index.html", "Nu s-a putut sterge")
		trans.Rollback()
		return
	}
	defer insertCustomer.Close()

	_, err = insertCustomer.Exec(session.Values["clientId"].(string))

	if err != nil {
		fmt.Println(err)
		temp.ExecuteTemplate(w, "index.html", "Nu s-a putut sterge")
		trans.Rollback()
		return
	}
	delete(session.Values, "clientId")
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "login.html", "Emailul a fost verificat cu succes!")
	trans.Commit()
	return

}
