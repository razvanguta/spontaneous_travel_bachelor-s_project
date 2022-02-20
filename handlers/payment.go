package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"myapp/database"
	"myapp/structs"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/jung-kurt/gofpdf"
	emailSend "github.com/scorredoira/email"
	"github.com/skip2/go-qrcode"
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
	sqlQueryA := "SELECT id,trip_title,trip_id from cart where client_id=?"
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
		var id, trip_title, trip_id string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&id, &trip_title, &trip_id) != nil {
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
		cartAll.TripID = trip_id
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

func GenerateQrPDF(percentage string, clientId string, text string, name string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 14)
	pdf.Cell(5, 5, text)
	var imgByte []byte
	imgByte, _ = qrcode.Encode(percentage, qrcode.Medium, 256)
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		return err
	}
	out, err := os.Create("assets\\qrImages\\" + percentage + "qr.jpg")
	if err != nil {
		return err
	}
	defer out.Close()
	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)

	if err != nil {
		return err
	}
	pdf.ImageOptions(
		"assets\\qrImages\\"+percentage+"qr.jpg",
		80, 20,
		0, 0,
		false,
		gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
		0,
		"",
	)
	err = pdf.OutputFileAndClose("assets\\pdf\\" + name + ".pdf")
	if err != nil {
		return err
	}

	var trans *sql.Tx
	trans, err = database.Db.Begin()
	if err != nil {
		return err
	}
	defer trans.Rollback()
	var insertDiscount *sql.Stmt
	insertDiscount, err = trans.Prepare("INSERT INTO discount (percentage, client_id) VALUES (?,?);")
	if err != nil {
		trans.Rollback()
		return err
	}
	defer insertDiscount.Close()
	_, err = insertDiscount.Exec(percentage, clientId)
	if err != nil {
		trans.Rollback()
		return err
	}

	err = trans.Commit()

	if err != nil {
		trans.Rollback()
		return err
	}

	return nil
}

func SendDiscountPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	temp.ExecuteTemplate(w, "sendDiscount.html", nil)
}

func SendDiscount(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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

		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var m structs.Comment
	if len(r.FormValue("description")) <= 0 || len(r.FormValue("username")) <= 0 || len(r.FormValue("percentage")) == 0 || len(r.FormValue("namePDF")) == 0 {
		m.ErrMessage = "Exista campuri goale!"
		temp.ExecuteTemplate(w, "sendDiscount.html", m)
		return
	}

	//take the client id
	sqlQueryA := "SELECT id, email from clients where username=?"
	rowA := database.Db.QueryRow(sqlQueryA, r.FormValue("username"))
	var id, email string
	//if don't exist => no agency
	if rowA.Scan(&id, &email) != nil {
		m.ErrMessage = "Nu exista acest client!"
		temp.ExecuteTemplate(w, "sendDiscount.html", m)
		return
	}

	err := GenerateQrPDF(r.FormValue("percentage"), id, r.FormValue("description"), r.FormValue("namePDF"))
	if err != nil {
		m.ErrMessage = "Ceva nu a mers cum trebuie1!"
		fmt.Println(err)
		temp.ExecuteTemplate(w, "sendDiscount.html", m)
		return
	}

	err = sendEmailQr(r.FormValue("username"), "Buna, "+r.FormValue("username")+", codul tau de reducere se afla mai jos!", "assets\\pdf\\"+r.FormValue("namePDF")+".pdf", email)
	if err != nil {
		m.ErrMessage = "Ceva nu a mers cum trebuie2!"
		fmt.Println(err)
		temp.ExecuteTemplate(w, "sendDiscount.html", m)
		return
	}

	m.ID = "yes"
	m.ErrMessage = "Cod trimis cu succes!"
	temp.ExecuteTemplate(w, "index.html", m)

}

func BuyTripPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	//take the client id
	sqlQueryA := "SELECT price, agencyID from trips where id=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("tripId"))
	var price, agencyID string
	//if don't exist => no agency
	if rowA.Scan(&price, &agencyID) != nil {
		message.ID = "yes"
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var trip structs.Trip
	trip.ID = param.ByName("tripId")
	trip.Price = price
	trip.ClientID = session.Values["Id"].(string)
	trip.AgencyID = agencyID
	temp.ExecuteTemplate(w, "buyPage.html", trip)
}

func BuyTripLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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

	name := r.FormValue("name")
	address := r.FormValue("address")
	phone := r.FormValue("phone")

	if len(name) == 0 || len(phone) == 0 || len(address) == 0 {
		message.ErrMessage = "Numarul de telefon contine alte caractere inafara de cifre!"
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	for _, c := range phone {
		if c < '0' || c > '9' {
			message.ErrMessage = "Numarul de telefon contine alte caractere inafara de cifre!"
			message.ID = "yes"
			temp.ExecuteTemplate(w, "index.html", message)
			return
		}
	}

	sqlQuery := "SELECT * from trips WHERE id=?"
	row := database.Db.QueryRow(sqlQuery, param.ByName("tripId"))
	var id, title, description, hotel, stars, price, img1, img2, img3, date, agencyID, numberOfDays, city, country string
	//if the username and path to profile is not the same redirect
	if row.Scan(&id, &title, &description, &hotel, &stars, &price, &img1, &img2, &img3, &date, &agencyID, &numberOfDays, &city, &country) != nil {
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

	//take the money to add
	sqlQueryA := "SELECT money_balance, email from clients where id=?"
	rowA := database.Db.QueryRow(sqlQueryA, session.Values["Id"].(string))
	var money_balance, email string
	//if don't exist => no agency
	if rowA.Scan(&money_balance, &email) != nil {
		var message structs.Comment
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	moneySum, err := strconv.Atoi(money_balance)
	if err != nil {
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	priceSum, err := strconv.Atoi(price)
	if err != nil {
		message.ErrMessage = "Ceva nu a mers cum trebuie!"
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	if priceSum > moneySum {
		message.ErrMessage = "Fond insuficient! Te rugam sa mai adaugi bani!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		return
	}

	var trans *sql.Tx
	trans, err = database.Db.Begin()
	if err != nil {
		message.ErrMessage = "Ceva nu a mers cum trebuie cu tranzactia!"
		message.ID = "yes"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var insertBought *sql.Stmt
	insertBought, err = trans.Prepare("INSERT INTO bought_trips (trip_title, trip_description, trip_hotel, trip_price, trip_numberOfDays, trip_city, trip_country, trip_date, agencyID, clientID, client_fullName, client_streetAddress, client_email, client_phone, path_to_pdf) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);")

	if err != nil {
		message.ID = "yes"
		fmt.Println(err)
		message.ErrMessage = "Ceva nu a mers cum trebuie cu inserarea1!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	defer insertBought.Close()

	_, err = insertBought.Exec(title, description, hotel, price, numberOfDays, city, country, date, param.ByName("agencyId"), session.Values["Id"].(string), r.FormValue("name"), r.FormValue("address"), email, r.FormValue("phone"), "assets\\pdf\\"+name+city+date+".pdf")

	if err != nil {
		message.ID = "yes"
		fmt.Println(err)
		message.ErrMessage = "Ceva nu a mers cum trebuie cu inserarea2!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	var deleteTrip *sql.Stmt

	deleteTrip, err = trans.Prepare("DELETE FROM trips where id=?")
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut sterge!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer deleteTrip.Close()
	_, err = deleteTrip.Exec(param.ByName("tripId"))

	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut sterge!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer deleteTrip.Close()

	var updateMoney *sql.Stmt
	updateMoney, err = trans.Prepare("UPDATE clients SET money_balance=? where id=?")
	if err != nil {
		fmt.Println(err)
		message.ErrMessage = "Ceva nu a mers cum trebuie cu suma de bani!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		trans.Rollback()
		return
	}
	defer updateMoney.Close()

	_, err = updateMoney.Exec(moneySum-priceSum, session.Values["Id"].(string))
	if err != nil {
		fmt.Println(err)
		message.ErrMessage = "Ceva nu a mers cum trebuie cu suma de bani!"
		temp.ExecuteTemplate(w, "addMoney.html", message)
		trans.Rollback()
		return
	}
	defer updateMoney.Close()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "I", 14)
	pdf.Cell(5, 5, "Salut "+name+", multumim pentru achizitie!\n")
	pdf.Ln(13)
	pdf.Cell(5, 5, "Aici se afla dovada platii excursiei in orasul "+city+" pe data de "+date+"un numar de "+numberOfDays+" zile\n")
	pdf.Ln(13)
	pdf.SetFont("Times", "B", 14)
	pdf.Cell(5, 5, "Pretul achizitiei "+price+" euro!\n")
	pdf.Ln(13)
	pdf.SetFont("Times", "I", 14)
	pdf.Cell(5, 5, "Un angajat al companiei te va contacta in legatura cu ultimele detalii!\n")
	pdf.Ln(13)
	err = pdf.OutputFileAndClose("assets\\pdf\\" + name + city + date + ".pdf")
	if err != nil {
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Ceva nu a mers bine cu factura!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}

	from := os.Getenv("EMAIL_ST")
	pass := os.Getenv("PASS_ST")
	m := emailSend.NewMessage("Achizitie excursie in "+city, "Buna multumim pentru achizitie, mai jos e dovada platii!")
	receiver := []string{email}
	m.From.Address = from
	m.From.Name = "admin spontaneous travel"
	m.To = receiver
	err = m.Attach("assets\\pdf\\" + name + city + date + ".pdf")

	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie cu factura!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	err = emailSend.Send("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), m)

	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie cu email!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	err = trans.Commit()

	if err != nil {
		message.ID = "yes"
		message.ErrMessage = "Ceva nu a mers cum trebuie cu tranzactia!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	message.ID = "yes"
	message.ErrMessage = "Excursie achizitionata cu succes!"
	temp.ExecuteTemplate(w, "index.html", message)
}

func JsonClientBoughtTrips(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	sqlQueryA := "SELECT id, trip_title, trip_description, trip_hotel, trip_city, trip_country, trip_date, client_fullName, path_to_pdf from bought_trips where clientID=?"
	rows, err := database.Db.Query(sqlQueryA, session.Values["Id"].(string))
	if err != nil {
		fmt.Println(err)
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
	trips := make([]*structs.Trip, 0)
	for rows.Next() {
		trip := new(structs.Trip)
		var id, trip_title, trip_description, trip_hotel, trip_city, trip_country, trip_date, client_fullName, path_to_pdf string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&id, &trip_title, &trip_description, &trip_hotel, &trip_city, &trip_country, &trip_date, &client_fullName, &path_to_pdf) != nil {
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
		trip.ID = id
		trip.ClientID = client_fullName
		trip.Title = trip_title
		trip.Description = trip_description
		trip.Hotel = trip_hotel
		trip.City = trip_city
		trip.Country = trip_country
		trip.Date = trip_date
		trip.PathToPDF = path_to_pdf
		trips = append(trips, trip)
	}
	// make an result of type Agency in order to create a json to send for the html page
	out, err := json.MarshalIndent(trips, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func ClientBoughtTrips(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	temp.ExecuteTemplate(w, "clientBoughtTrips.html", message)
}

func JsonAgencyBoughtTrips(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	if session.Values["Role"] != "AGENCY" {
		message.ID = "yes"
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	sqlQueryA := "SELECT id,trip_hotel, trip_city, trip_date, client_fullName,client_streetAddress,client_email,client_phone,path_to_pdf from bought_trips where agencyID=?"
	rows, err := database.Db.Query(sqlQueryA, session.Values["Id"].(string))
	if err != nil {
		fmt.Println(err)
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
	clients := make([]*structs.Client, 0)
	for rows.Next() {
		client := new(structs.Client)
		var id, trip_hotel, trip_city, trip_date, client_fullName, client_streetAddress, client_email, client_phone, path_to_pdf string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&id, &trip_hotel, &trip_city, &trip_date, &client_fullName, &client_streetAddress, &client_email, &client_phone, &path_to_pdf) != nil {
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
		client.ID = id
		client.Name = client_fullName
		client.Phone = client_phone
		client.Address = client_streetAddress
		client.PathToPDF = path_to_pdf
		client.City = trip_city
		client.Date = trip_date
		client.Hotel = trip_hotel
		clients = append(clients, client)
	}
	// make an result of type Agency in order to create a json to send for the html page
	out, err := json.MarshalIndent(clients, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func AgencyBoughtTrips(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	if session.Values["Role"] != "AGENCY" {
		message.ID = "yes"
		message.ErrMessage = "Nu poti efectua aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	temp.ExecuteTemplate(w, "agencyBoughtTrips.html", message)
}
