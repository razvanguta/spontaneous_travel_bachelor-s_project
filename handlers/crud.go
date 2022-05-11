package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func DeleteAgencyPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		var message structs.Comment
		message.Username = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if session.Values["Role"] != "ADMIN" {
		var message structs.Comment
		message.Username = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	temp.ExecuteTemplate(w, "deleteClient.html", nil)
}

func DeleteClient(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		var message structs.Comment
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if session.Values["Role"] != "ADMIN" {
		var message structs.Comment
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var trans *sql.Tx
	trans, err := database.Db.Begin()
	if err != nil {
		var message structs.Comment
		message.Message = "A aparut o eroare, te rog mai incearca!"
		temp.ExecuteTemplate(w, "register.html", message)
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()
	var deleteAgency *sql.Stmt
	deleteAgency, err = trans.Prepare("DELETE FROM clients where username=?")
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu s-a putut sterge!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer deleteAgency.Close()
	_, err = deleteAgency.Exec(r.FormValue("username"))
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu s-a putut sterge!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer deleteAgency.Close()
	trans.Commit()
	var message structs.Comment
	message.ID = "yes"
	message.ErrMessage = "Ai sters cu succes clientul!"
	temp.ExecuteTemplate(w, "index.html", message)
	trans.Rollback()
}

func DeleteAgency(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		var message structs.Comment
		message.Username = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if session.Values["Role"] != "ADMIN" {
		session.Options.MaxAge = -1
		session.Save(r, w)
		var message structs.Comment
		message.Username = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var trans *sql.Tx
	trans, err := database.Db.Begin()
	if err != nil {
		var message structs.Comment
		message.Message = "A aparut o eroare, te rog mai incearca!"
		temp.ExecuteTemplate(w, "register.html", message)
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()
	var deleteAgency *sql.Stmt
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
	_, err = deleteAgency.Exec(param.ByName("agencyId"))
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
	trans.Commit()
	http.Redirect(w, r, "/allAgencies", 301)
}

func DeleteMyself(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		var message structs.Comment
		message.Username = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	var trans *sql.Tx
	trans, err := database.Db.Begin()

	if err != nil {
		var message structs.Comment
		message.Message = "A aparut o eroare, te rog mai incearca!"
		temp.ExecuteTemplate(w, "register.html", message)
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
		temp.ExecuteTemplate(w, "editDescriptionAgency.html", "A aparut o eroare, te rog mai incearca!")
		return
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
		temp.ExecuteTemplate(w, "editDescriptionAgency.html", "A aparut o eroare, te rog mai incearca!")
		return
	}
	fmt.Println(r)
	if len(r.FormValue("description")) < 25 {
		temp.ExecuteTemplate(w, "editDescriptionAgency.html", "Introdu o descriere care contine mai mult de 25 de caractere!")
		return
	}
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
	message.ID = "yes"
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
		var m structs.Comment
		m.ErrMessage = "Parola nu respecta criteriile"
		temp.ExecuteTemplate(w, "changePassword.html", m)
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
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Esti conectat!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	temp.ExecuteTemplate(w, "passwordReset.html", nil)
}

func ReviewJson(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//select in another query the particularities from reviews
	sqlQueryA := "SELECT client_id,agency_id,title,comment,stars,date from reviews where agency_id=?"
	rows, err := database.Db.Query(sqlQueryA, param.ByName("agencyId"))
	if err != nil {
		session, _ := Store.Get(r, "session")
		//we send a message if the user is connected so that some buttons will not display
		var message structs.Comment
		if !session.IsNew {
			message.ErrMessage = "Ceva nu a mers cum trebuie!"
			message.ID = "yes"
			temp.ExecuteTemplate(w, "index.html", message)
		} else {
			session.Options.MaxAge = -1
			session.Save(r, w)
			message.ErrMessage = "Ceva nu a mers cum trebuie!"
			temp.ExecuteTemplate(w, "index.html", message)
		}
		return
	}
	reviews := make([]*structs.Review, 0)
	for rows.Next() {
		review := new(structs.Review)
		var client_id, agency_id, title, comment, stars, date string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&client_id, &agency_id, &title, &comment, &stars, &date) != nil {
			session, _ := Store.Get(r, "session")
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
		sqlQueryA := "SELECT username from clients where id=?"
		rowA := database.Db.QueryRow(sqlQueryA, client_id)
		var username string
		//if don't exist => no agency
		if rowA.Scan(&username) != nil {
			var message structs.Comment
			//we send a message if the user is connected so that some buttons will not display
			message.ID = "yes"
			message.ErrMessage = "Ceva nu a mers cum trebuie!"
			temp.ExecuteTemplate(w, "index.html", message)
			return
		}
		review.Client = username
		review.Title = title
		review.Comment = comment
		review.Stars = stars
		review.Date = date
		session, _ := Store.Get(r, "session")
		if !session.IsNew {
			if session.Values["Role"].(string) == "ADMIN" {
				review.IsTheSame = "yes"
			} else if session.Values["Id"] == client_id {
				review.IsTheSame = "yes"
			} else {
				review.IsTheSame = "no"
			}
		} else {
			session.Options.MaxAge = -1
			session.Save(r, w)
			review.IsTheSame = "no"
		}
		reviews = append(reviews, review)
	}
	// make an result of type Agency in order to create a json to send for the html page
	out, err := json.MarshalIndent(reviews, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func PutReview(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	var message structs.Comment
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if session.Values["Role"] != "CLIENT" {
		message.ID = "yes"
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if len(r.FormValue("title")) == 0 || len(r.FormValue("description")) == 0 {
		message.ID = "yes"
		message.ErrMessage = "Nu poti lasa campuri necompletate!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if !session.IsNew {
		if session.Values["Role"] == "CLIENT" {
			message.IsClient = "yes"
			// start the transaction, because all the validations passed at this point
			var trans *sql.Tx
			trans, err := database.Db.Begin()
			if err != nil {
				temp.ExecuteTemplate(w, "index.html", "A aparut o eroare, te rog mai incearca!")
				return
			}
			//this will be ignored in case of a commit
			defer trans.Rollback()
			//insert review
			var insertReview *sql.Stmt
			insertReview, err = trans.Prepare("INSERT INTO reviews (client_id, agency_id, title, comment, stars, date) VALUES (?,?,?,?,?,?);")
			if err != nil {
				message.ErrMessage = "Nu s-a putut adauga"
				temp.ExecuteTemplate(w, "index.html", message)
				fmt.Println(err)
				trans.Rollback()
				return
			}
			defer insertReview.Close()
			_, err = insertReview.Exec(param.ByName("userId"), param.ByName("agencyId"), r.FormValue("title"), r.FormValue("description"), r.FormValue("points"), time.Now().Format("01-02-2006 15:04:05"))
			if err != nil {
				message.ErrMessage = "Ceva nu a functionat cum trebuie!"
				temp.ExecuteTemplate(w, "index.html", message)
				fmt.Println(err)
				trans.Rollback()
				return
			}

			err = trans.Commit()

			if err != nil {
				message.ErrMessage = "Ceva nu a functionat cum trebuie!"
				temp.ExecuteTemplate(w, "index.html", message)
				trans.Rollback()
				return
			}
		} else {
			message.ID = "yes"
			message.ErrMessage = "Nu poti accesa!"
			temp.ExecuteTemplate(w, "index.html", message)
			return
		}
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
		message.ErrMessage = "Nu poti accesa!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	sqlQueryA := "SELECT username from agencies where id=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("agencyId"))
	var username string
	//if don't exist => no agency
	if rowA.Scan(&username) != nil {
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

	http.Redirect(w, r, "/seeReviews/"+username, 301)
}

func SeeReviews(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	var message structs.Comment
	if !session.IsNew {
		message.ID = "yes"
		if session.Values["Role"] == "CLIENT" {
			message.ID = session.Values["Id"].(string)
			message.IsClient = "yes"
		}
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	sqlQueryA := "SELECT id from agencies where username=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("nameOfAgency"))
	var id string
	//if don't exist => no agency
	if rowA.Scan(&id) != nil {
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
	message.ID2 = id
	temp.ExecuteTemplate(w, "reviewPage.html", message)
}

func DeleteReview(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
		var message structs.Comment
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	fmt.Println(1)
	if session.Values["Role"] != "CLIENT" && session.Values["Role"].(string) != "ADMIN" {
		session, _ := Store.Get(r, "session")
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
	fmt.Println(2)
	if session.Values["Role"].(string) == "CLIENT" && session.Values["Id"].(string) != param.ByName("clientId") {
		session, _ := Store.Get(r, "session")
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
	fmt.Println(3)
	var trans *sql.Tx
	trans, err := database.Db.Begin()

	if err != nil {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Ceva nu este bine!"
		temp.ExecuteTemplate(w, "index.html", message)
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	var deleteReview *sql.Stmt

	deleteReview, err = trans.Prepare("DELETE FROM reviews where agency_id=? and date=?")
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut sterge!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer deleteReview.Close()
	_, err = deleteReview.Exec(param.ByName("agencyId"), param.ByName("date"))

	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut sterge!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer deleteReview.Close()
	trans.Commit()

	sqlQueryA := "SELECT username from agencies where id=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("agencyId"))
	var username string
	//if don't exist => no agency
	if rowA.Scan(&username) != nil {
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
	http.Redirect(w, r, "/seeReviews/"+username, 301)

}

func PasswordResetLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Esti conectat!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	email := r.FormValue("email")

	err := checkEmail(email)

	if err != nil {
		var m structs.Comment
		m.ErrMessage = "Email invalid!"
		temp.ExecuteTemplate(w, "passwordReset.html", m)
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
		temp.ExecuteTemplate(w, "passwordReset.html", nil)
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
	fmt.Println(isAgency)
	sqlQuerry2 := "SELECT ROLE FROM CLIENTS WHERE EMAIL=?"
	row2 := database.Db.QueryRow(sqlQuerry2, email)
	var role2 string

	if !isAgency && row2.Scan(&role2) != nil {
		temp.ExecuteTemplate(w, "passwordReset.html", "Email-ul nu exista")
		return
	}
	//hash the password in order to crypt
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		temp.ExecuteTemplate(w, "passwordReset.html", "Nu s-a putut inregistra1")
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
			temp.ExecuteTemplate(w, "passwordReset.html", message)
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
			temp.ExecuteTemplate(w, "passwordReset.html", message)
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
			temp.ExecuteTemplate(w, "passwordReset.html", message)
			temp.ExecuteTemplate(w, "passwordReset.html", message)
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
			temp.ExecuteTemplate(w, "passwordReset.html", message)
			temp.ExecuteTemplate(w, "passwordReset.html", message)
			trans.Rollback()
			return
		}
		defer updateCustomerPassword.Close()
	}
	var message structs.Comment
	message.ID = "yes"
	message.ErrMessage = "Parola resetata cu succes!"
	temp.ExecuteTemplate(w, "passwordReset.html", message)
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
