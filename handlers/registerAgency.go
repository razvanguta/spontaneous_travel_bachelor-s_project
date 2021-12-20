package handlers

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"html/template"
	"math/big"
	"myapp/database"
	"myapp/structs"
	"net/http"
	"net/smtp"
	"os"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAgency(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if !session.IsNew && session.Values["Role"].(string) != "ADMIN" {
		var message structs.Comment
		if session.Values["Role"].(string) == "AGENCY" {
			message.IsAgency = "yes"
			message.Username = session.Values["Username"].(string)
			message.ErrMessage = "Esti deja conectat!"
			message.IsTheSame = "yes"
			temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
			return
		}
		message.Username = "Esti deja conectat!"
		message.ID = "yes"
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		temp.ExecuteTemplate(w, "personalPage.html", message)
		return
	}
	if session.Values["Role"].(string) != "ADMIN" {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	temp.ExecuteTemplate(w, "registerAgency.html", nil)
}

func RegisterAgencyLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if !session.IsNew && session.Values["Role"].(string) != "ADMIN" {
		var message structs.Comment
		if session.Values["Role"].(string) == "AGENCY" {
			message.IsAgency = "yes"
			message.Username = session.Values["Username"].(string)
			message.ErrMessage = "Esti deja conectat!"
			message.IsTheSame = "yes"
			temp.ExecuteTemplate(w, "agencyPersonalPage.html", message)
			return
		}
		message.Username = "Esti deja conectat!"
		message.ID = "yes"
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		temp.ExecuteTemplate(w, "personalPage.html", message)
		return
	}
	if session.Values["Role"].(string) != "ADMIN" {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	err1 := r.ParseForm()

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	//take username from form and check to have the correct criteria
	var username string = r.FormValue("username")

	err := checkUsername(username)

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", err.Error())
		return
	}

	//generate random secured password that will be send by email
	var password string = generateRandomPassword()
	for checkPassword(password) != nil {
		password = generateRandomPassword()
	}
	//check email
	var email string = r.FormValue("email")

	err = checkEmail(email)

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", err.Error())
		return
	}

	//TODO: PATH TO PHOTO

	// start the transaction, because all the validations passed at this point
	var trans *sql.Tx
	trans, err = database.Db.Begin()
	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", "A aparut o eroare, te rog mai incearca!")
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	sqlQuerry := "SELECT id FROM agencies WHERE username=?;"
	row := trans.QueryRow(sqlQuerry, username)
	var id string
	if row.Scan(&id) != sql.ErrNoRows {
		temp.ExecuteTemplate(w, "registerAgency.html", "Ne pare rau dar usernameul nu este valid")
		trans.Rollback()
		return
	}

	//hash the password in order to crypt
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		temp.ExecuteTemplate(w, "registerAgency.html", "Nu s-a putut inregistra1")
		trans.Rollback()
		return
	}

	//insert the agency
	var insertAgency *sql.Stmt
	insertAgency, err = trans.Prepare("INSERT INTO agencies (username, description, email, hash, role, is_active,path_profile_image) VALUES (?,?,?,?,?,?,?);")

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", "Nu s-a putut inregistra2")
		fmt.Println(err)
		trans.Rollback()
		return
	}
	defer insertAgency.Close()

	_, err = insertAgency.Exec(username, "Te rugam sa adaugi descrierea agentiei tale aici", email, hash, "AGENCY", 1, "firma1.jpg")

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", "Nu s-a putut inregistra3")
		fmt.Println(err)
		trans.Rollback()
		return
	}
	err = sendUsernamePassword(username, password, email)

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", "Nu s-a putut inregistra6")
		fmt.Println(err)
		trans.Rollback()
		return
	}

	err = trans.Commit()

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", "Nu s-a putut inregistra7")
		trans.Rollback()
		return
	}
	var m structs.Comment

	m.Username = "Agentia a fost adaugata cu succes"
	m.ID = "yes"
	temp.ExecuteTemplate(w, "index.html", m)

}

func sendUsernamePassword(username string, password string, receiverEmail string) error {
	from := os.Getenv("EMAIL_ST")
	pass := os.Getenv("PASS_ST")
	receiver := []string{receiverEmail}
	address := "smtp.gmail.com:587"

	subject := "Subject: Nume si Parola pentru contul tau de agent de turism\r\n\r\n"

	body := "Buna, multumim ca te-ai alaturat noua, numele de utilizator este: " + username + " si parola: " + password

	mail := []byte(subject + body)

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	err := smtp.SendMail(address, auth, from, receiver, mail)

	return err

}

func generateRandomPassword() string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!?#"
	pass := make([]byte, 8)
	for i := 0; i < 8; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			panic(err)
		}
		pass[i] = alphabet[num.Int64()]
	}
	return string(pass)
}
