package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"
	"net/smtp"
	"os"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
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

func ChangePassword(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", nil)
		return
	}
	temp.ExecuteTemplate(w, "changePassword.html", nil)
}

func ChangePasswordLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		temp.ExecuteTemplate(w, "index.html", nil)
		return
	}

	var trans *sql.Tx
	trans, err := database.Db.Begin()

	if err != nil {
		temp.ExecuteTemplate(w, "index.html", nil)
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	var updateCustomerPassword *sql.Stmt
	var updateAgencyPassword *sql.Stmt

	password := r.FormValue("password")
	err = checkPassword(password)

	if err != nil {
		temp.ExecuteTemplate(w, "changePassword.html", "Parola nu respecta criteriile")
		return
	}
	//hash the password in order to crypt
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		temp.ExecuteTemplate(w, "index.html", "Nu s-a putut inregistra1")
		trans.Rollback()
		return
	}

	if session.Values["Role"].(string) == "AGENCY" {
		updateAgencyPassword, err = trans.Prepare("UPDATE agencies SET hash=? where id=?")
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut edita!"
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer updateAgencyPassword.Close()
		_, err = updateAgencyPassword.Exec(hash, session.Values["Id"].(string))
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut edita!"
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer updateAgencyPassword.Close()

	} else {
		updateCustomerPassword, err = trans.Prepare("UPDATE clients SET hash=? where id=?")
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut edita!"
			temp.ExecuteTemplate(w, "index.html", message)
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer updateCustomerPassword.Close()
		_, err = updateCustomerPassword.Exec(hash, session.Values["Id"].(string))
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut edita!"
			temp.ExecuteTemplate(w, "index.html", message)
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer updateCustomerPassword.Close()
	}
	var message structs.Comment
	message.ID = "yes"
	message.ErrMessage = "Parola schimbata cu succes!"
	temp.ExecuteTemplate(w, "index.html", message)
	trans.Commit()

}

func PasswordReset(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	temp.ExecuteTemplate(w, "passwordReset.html", nil)
}

func PasswordResetLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	email := r.FormValue("email")

	err := checkEmail(email)

	if err != nil {
		temp.ExecuteTemplate(w, "passwordReset.html", "Email invalid!")
		return
	}

	//generate random secured password that will be send by email
	var password string = generateRandomPassword()
	for checkPassword(password) != nil {
		password = generateRandomPassword()
	}
	var trans *sql.Tx
	trans, err = database.Db.Begin()
	if err != nil {
		temp.ExecuteTemplate(w, "index.html", nil)
		return
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	var updateAgencyPassword *sql.Stmt
	var updateCustomerPassword *sql.Stmt
	var isAgency bool

	sqlQuerry := "SELECT ROLE FROM AGENCIES WHERE EMAIL=?"
	row := database.Db.QueryRow(sqlQuerry, email)
	var role string
	isAgency = true
	if row.Scan(&role) != nil {
		isAgency = false
	}
	sqlQuerry2 := "SELECT ROLE FROM CLIENTS WHERE EMAIL=?"
	row2 := database.Db.QueryRow(sqlQuerry2, email)
	var role2 string

	if row2.Scan(&role2) != nil {
		temp.ExecuteTemplate(w, "passwordReset.html", "Email-ul nu exista")
		return
	}
	//hash the password in order to crypt
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		temp.ExecuteTemplate(w, "index.html", "Nu s-a putut inregistra1")
		trans.Rollback()
		return
	}

	if isAgency {
		updateAgencyPassword, err = trans.Prepare("UPDATE agencies SET hash=? where email=?")
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut sterge!"
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer updateAgencyPassword.Close()
		_, err = updateAgencyPassword.Exec(hash, email)
		if err != nil {
			fmt.Println(err)
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Nu s-a putut sterge!"
			temp.ExecuteTemplate(w, "index.html", message)
			trans.Rollback()
			return
		}
		defer updateAgencyPassword.Close()

	} else {
		updateCustomerPassword, err = trans.Prepare("UPDATE clients SET hash=? where email=?")
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
		defer updateCustomerPassword.Close()
		_, err = updateCustomerPassword.Exec(hash, email)
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
		defer updateCustomerPassword.Close()
	}
	var message structs.Comment
	message.ID = "yes"
	message.ErrMessage = "Parola resetata cu succes!"
	temp.ExecuteTemplate(w, "index.html", message)
	sendPassword(password, email)
	trans.Commit()

}

func sendPassword(password string, receiverEmail string) error {
	from := os.Getenv("EMAIL_ST")
	pass := os.Getenv("PASS_ST")
	receiver := []string{receiverEmail}
	address := "smtp.gmail.com:587"

	subject := "Subject: Parola dupa resetare\r\n\r\n"

	body := "Buna, noua ta parola este: " + password + " , aceasta se poate schimba oricand!"

	mail := []byte(subject + body)

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	err := smtp.SendMail(address, auth, from, receiver, mail)

	return err

}
