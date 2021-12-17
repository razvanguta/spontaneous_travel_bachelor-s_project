package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"myapp/database"
	"myapp/structs"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var temp *template.Template

func RegisterClientPage(w http.ResponseWriter, r *http.Request) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		var message structs.Comment
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		message.Username = "Esti deja conectat!"
		temp.ExecuteTemplate(w, "personalPage.html", message)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "register.html", nil)
}

func RegisterClientLogic(w http.ResponseWriter, r *http.Request) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		var message structs.Comment
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		message.Username = "Esti deja conectat"
		temp.ExecuteTemplate(w, "personalPage.html", message)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	err1 := r.ParseForm()

	if err1 != nil {
		fmt.Println(err1)
		return
	}
	//take username from form and check to have the correct criteria
	var username string = r.FormValue("username")

	err := checkUsername(username)

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", err.Error())
		return
	}

	//same for password
	var password string = r.FormValue("password")

	err = checkPassword(password)

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", err.Error())
		return
	}

	//same for email
	var email string = r.FormValue("email")

	err = checkEmail(email)

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", err.Error())
		return
	}

	// start the transaction, because all the validations passed at this point
	var trans *sql.Tx
	trans, err = database.Db.Begin()
	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "A aparut o eroare, te rog mai incearca!")
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	sqlQuerry := "SELECT id FROM clients WHERE username=?;"
	row := trans.QueryRow(sqlQuerry, username)
	var id string
	if row.Scan(&id) != sql.ErrNoRows {
		temp.ExecuteTemplate(w, "register.html", "Ne pare rau dar usernameul nu este valid")
		trans.Rollback()
		return
	}

	//hash the password in order to crypt
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra1")
		trans.Rollback()
		return
	}
	//insert the user
	var insertCustomer *sql.Stmt
	insertCustomer, err = trans.Prepare("INSERT INTO clients (username, money_balance, email, is_active, hash, role) VALUES (?,?,?,?,?,?);")

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra2")
		trans.Rollback()
		return
	}
	defer insertCustomer.Close()

	_, err = insertCustomer.Exec(username, 0, email, 0, hash, "CLIENT")

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra3")
		fmt.Println(err)
		trans.Rollback()
		return
	}

	//create random code in order to verify the email and insert
	rand.Seed(time.Now().UnixNano())

	verCode := rand.Intn(99999999)
	var insertEmailVer *sql.Stmt
	insertEmailVer, err = trans.Prepare("INSERT INTO email_ver (username, email, verification_code) VALUES (?, ?, ?);")

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra4")
		fmt.Println(err)
		trans.Rollback()
		return
	}

	defer insertEmailVer.Close()

	_, err = insertEmailVer.Exec(username, email, verCode)

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra5")
		fmt.Println(err)
		trans.Rollback()
		return
	}
	err = sendVerCode(verCode, email)

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra6")
		fmt.Println(err)
		trans.Rollback()
		return
	}

	err = trans.Commit()

	if err != nil {
		temp.ExecuteTemplate(w, "register.html", "Nu s-a putut inregistra7")
		trans.Rollback()
		return
	}

	var m structs.Comment

	m.Email = email
	temp.ExecuteTemplate(w, "emailVerification.html", m)
}

func sendVerCode(verCode int, receiverEmail string) error {
	from := os.Getenv("EMAIL_ST")
	pass := os.Getenv("PASS_ST")
	receiver := []string{receiverEmail}
	address := "smtp.gmail.com:587"

	subject := "Subject: Cod de verificare pentru Spontaneous Travel\r\n\r\n"

	verCodeStr := strconv.Itoa(verCode)

	body := "Buna, multumim ca te-ai alaturat noua, codul de verificare este: " + verCodeStr

	mail := []byte(subject + body)

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	err := smtp.SendMail(address, auth, from, receiver, mail)

	return err

}

func EmailVerification(w http.ResponseWriter, r *http.Request) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		var message structs.Comment
		if session.Values["Role"].(string) == "ADMIN" {
			message.IsAdmin = "yes"
		}
		temp.ExecuteTemplate(w, "personalPage.html", "Esti deja conectat")
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	err1 := r.ParseForm()

	if err1 != nil {
		fmt.Println(err1)
		return
	}
	//extrect what user inserted for email and verification code
	email := r.FormValue("email")
	verCode := r.FormValue("vercode")

	trans, err := database.Db.Begin()
	if err != nil {
		temp.ExecuteTemplate(w, "emailVerification.html", "A aparut o eroare, te rog mai incearca!")
	}

	defer trans.Rollback()

	sqlQuery := "SELECT verification_code FROM email_ver WHERE email=?;"
	row := trans.QueryRow(sqlQuery, email)
	var verification_code string
	row.Scan(&verification_code)
	// if the ver code is the same with the verification code sent in the email we update
	if verCode == verification_code {
		sqlQuery2 := "UPDATE clients SET is_active = 1 WHERE email = ?;"
		update, err := trans.Prepare(sqlQuery2)

		if err != nil {
			temp.ExecuteTemplate(w, "emailVerification.html", "Nu am putut verifica, mai incercati!")
			return
		}
		defer update.Close()

		_, err = update.Exec(email)

		if err != nil {
			fmt.Println(err)
			temp.ExecuteTemplate(w, "emailVerification.html", "Nu s-a putut inregistra8")
			trans.Rollback()
			return
		}

		temp.ExecuteTemplate(w, "login.html", "Emailul a fost verificat cu succes!")
		trans.Commit()
		return

	}
	// otherwise we put the user to introduce the email one more time
	trans.Rollback()
	temp.ExecuteTemplate(w, "emailVerification.html", "A aparut o eroare")

}

func checkUsername(username string) error {
	for _, ch := range username {
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) {
			return errors.New("usernam-ul trebuie sa contina doar litere si cifre")
		}
	}
	return nil
}

func checkPassword(password string) error {
	var isLower, isUpper, isNumber bool
	isLower = false
	isUpper = false
	isNumber = false
	var err error
	for _, ch := range password {
		if unicode.IsLower(ch) {
			isLower = true
		}
		if unicode.IsUpper(ch) {
			isUpper = true
		}
		if unicode.IsNumber(ch) {
			isNumber = true
		}
		if unicode.IsSpace(ch) {
			err = errors.New("parola contine spatii")
		}
	}

	if !isLower {
		err = errors.New("parola trebuie sa contina litere mici")
	}
	if !isUpper {
		err = errors.New("parola trebuie sa contina litere mari")
	}
	if !isNumber {
		err = errors.New("parola trebuie sa contina cifre")
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

func checkEmail(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return errors.New("te rugam sa introduci o adresa de email valida")
	}

	a := strings.Index(email, "@")
	domain := email[a+1:]
	_, err = net.LookupMX(domain)

	if err != nil {
		return errors.New("nu exista acest domeniu de email")
	}

	return nil
}
