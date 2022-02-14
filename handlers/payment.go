package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AddToCart(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
	//we check before if we are connected, so this page will not display
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
	//take the title
	sqlQueryA := "SELECT title from trips where id=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("tripId"))
	var title string
	//if don't exist => no title
	if rowA.Scan(&title) != nil {
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
	//check if the trip is already in the cart
	sqlQuery := "SELECT trip_id,client_id,trip_title from cart"
	rows, err := database.Db.Query(sqlQuery)
	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	for rows.Next() {
		var trip_id, client_id, trip_title string
		if rows.Scan(&trip_id, &client_id, &trip_title) != nil {
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
		if trip_id == param.ByName("tripId") && client_id == param.ByName("clientId") && trip_title == title {
			message.ID = "yes"
			message.ErrMessage = "Exista deja in cos!"
			temp.ExecuteTemplate(w, "index.html", message)
			return
		}

	}

	var trans *sql.Tx
	trans, err = database.Db.Begin()
	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	var insertCart *sql.Stmt
	insertCart, err = trans.Prepare("INSERT INTO cart (trip_id, client_id, trip_title) VALUES (?,?,?);")

	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	defer insertCart.Close()

	fmt.Println(param.ByName("clientId"))
	_, err = insertCart.Exec(param.ByName("tripId"), param.ByName("clientId"), title)

	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	err = trans.Commit()

	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	http.Redirect(w, r, "/seeCart", 301)
}

func JsonSeeCart(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
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
	sqlQueryA := "SELECT id,trip_title from cart where client_id=?"
	rows, err := database.Db.Query(sqlQueryA, session.Values["Id"].(string))
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
	cart := make([]*structs.Trip, 0)
	for rows.Next() {
		cartAll := new(structs.Trip)
		var id, trip_title string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&id, &trip_title) != nil {
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
		cartAll.ID = id
		cartAll.Title = trip_title
		cart = append(cart, cartAll)
	}
	// make an result of type Agency in order to create a json to send for the html page
	out, err := json.MarshalIndent(cart, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func SeeCart(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
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
	temp.ExecuteTemplate(w, "seeCart.html", message)
}

func OutFromCart(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
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

	deleteReview, err = trans.Prepare("DELETE FROM cart where id=?")
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
	_, err = deleteReview.Exec(param.ByName("cartId"))

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
	http.Redirect(w, r, "/seeCart", 301)
}

func AddMoney(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
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
	temp.ExecuteTemplate(w, "addMoney.html", message)
}

func AddMoneyLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ := template.ParseGlob("templates/*.html")
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

	card := r.FormValue("card")
	cvv := r.FormValue("cvv")
	money := r.FormValue("money")

	if len(card) == 0 || len(cvv) == 0 || len(money) == 0 {
		message.ErrMessage = "Nu pot exista campuri goale!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		return
	}

	for _, c := range card {
		if c < '0' || c > '9' {
			message.ErrMessage = "Numarul cardului contine alte caractere inafara de cifre!"
			temp.ExecuteTemplate(w, "addMoney.html", message)
			return
		}
	}

	for _, c := range cvv {
		if c < '0' || c > '9' {
			message.ErrMessage = "CVV-ul contine alte caractere inafara de cifre!"
			temp.ExecuteTemplate(w, "addMoney.html", message)
			return
		}
	}

	//start the transaction
	var trans *sql.Tx
	trans, err := database.Db.Begin()

	if err != nil {
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		return
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	//take the money to add
	sqlQueryA := "SELECT money_balance from clients where id=?"
	rowA := database.Db.QueryRow(sqlQueryA, session.Values["Id"].(string))
	var money_balance string
	//if don't exist => no agency
	if rowA.Scan(&money_balance) != nil {
		var message structs.Comment
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		return
	}
	moneySum, err := strconv.Atoi(money)
	if err != nil {
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		return
	}
	moneyBalanceSum, err := strconv.Atoi(money_balance)
	if err != nil {
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		return
	}
	moneyS := (moneySum + moneyBalanceSum)
	//we start the porccess of description update
	var updateMoney *sql.Stmt
	updateMoney, err = trans.Prepare("UPDATE clients SET money_balance=? where id=?")
	if err != nil {
		fmt.Println(err)
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		trans.Rollback()
		return
	}
	defer updateMoney.Close()

	_, err = updateMoney.Exec(strconv.Itoa(moneyS), session.Values["Id"].(string))
	if err != nil {
		fmt.Println(err)
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		trans.Rollback()
		return
	}
	defer updateMoney.Close()
	trans.Commit()
	http.Redirect(w, r, "/myPersonalPage", 301)
}
