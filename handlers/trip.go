package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"myapp/database"
	"myapp/structs"
	"net/http"
	"os"
	"strconv"

	weather "github.com/briandowns/openweathermap"
	"github.com/julienschmidt/httprouter"
)

func DeleteTrip(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

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
	sqlQuery := "SELECT agencyID from trips WHERE id=?"
	row := database.Db.QueryRow(sqlQuery, param.ByName("tripId"))
	var agencyID string
	if row.Scan(&agencyID) != nil {
		fmt.Println(param.ByName("agencyID"))
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

	if session.Values["Role"].(string) != "AGENCY" && session.Values["Role"].(string) != "ADMIN" {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if session.Values["Role"].(string) == "AGENCY" && agencyID != session.Values["Id"] {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	//we check before if we are connected, so this page will not display
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
	var message structs.Comment
	message.ID = "yes"
	message.Username = "Excursie stearsa cu succes!"
	temp.ExecuteTemplate(w, "index.html", message)
	trans.Commit()
}

func JsonUpdateTrip(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	sqlQuery := "SELECT agencyID from trips WHERE id=?"
	row := database.Db.QueryRow(sqlQuery, param.ByName("tripId"))
	var agencyID string
	if row.Scan(&agencyID) != nil {
		fmt.Println(param.ByName("agencyID"))
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
	if session.Values["Role"].(string) != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if agencyID != session.Values["Id"] {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	sqlQueryA := "SELECT * from trips WHERE id=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("tripId"))
	var id, title, description, hotel, stars, price, img1, img2, img3, date, numberOfDays, city, country string
	//if the username and path to profile is not the same redirect
	if rowA.Scan(&id, &title, &description, &hotel, &stars, &price, &img1, &img2, &img3, &date, &agencyID, &numberOfDays, &city, &country) != nil {
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
	var trip structs.Trip
	trip.ID = id
	trip.Title = title
	trip.Description = description
	trip.Date = date
	trip.NumberOfDays = numberOfDays
	trip.Hotel = hotel
	trip.Stars = stars
	trip.Price = price
	trip.Path_img1 = img1
	trip.Path_img2 = img2
	trip.Path_img3 = img3
	trip.City = city
	trip.Country = country
	// make an result of type Trip in order to create a json to send for the html page
	out, err := json.MarshalIndent(trip, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func UpdateTripPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	sqlQuery := "SELECT agencyID from trips WHERE id=?"
	row := database.Db.QueryRow(sqlQuery, param.ByName("tripId"))
	var agencyID string
	if row.Scan(&agencyID) != nil {
		fmt.Println(param.ByName("agencyID"))
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
	if session.Values["Role"].(string) != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if agencyID != session.Values["Id"] {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var message structs.Comment
	message.ID = param.ByName("tripId")
	temp.ExecuteTemplate(w, "updateTrip.html", message)
}

func UpdateTripLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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

	sqlQuery := "SELECT agencyID from trips WHERE id=?"
	row := database.Db.QueryRow(sqlQuery, param.ByName("tripId"))
	fmt.Println(param.ByName("tripId"))
	var agencyID string
	if row.Scan(&agencyID) != nil {
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
	if session.Values["Role"].(string) != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	if agencyID != session.Values["Id"] {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu poti accesa aceasta operatiune!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var title string = r.FormValue("title")
	var description string = r.FormValue("description")
	var hotel string = r.FormValue("hotel")
	var stars string = r.FormValue("stars")
	var price string = r.FormValue("price")
	if priceFloat, err := strconv.ParseFloat(price, 64); err == nil {
		if priceFloat < 0 {
			var message structs.Comment
			message.ID = "yes"
			message.Username = "Pretul nu poate fi negativ!"
			temp.ExecuteTemplate(w, "index.html", message)
			return
		}
	}
	var date string = r.FormValue("date")
	var days string = r.FormValue("days")
	var country string = r.FormValue("country")
	var city string = r.FormValue("city")
	if len(title) <= 0 || len(description) <= 0 || len(hotel) <= 0 || len(stars) <= 0 || len(price) <= 0 || len(date) <= 0 || len(days) <= 0 || len(country) <= 0 || len(city) <= 0 {
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Au existat campuri goale!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	//we check before if we are connected, so this page will not display
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

	var updateTrip *sql.Stmt

	updateTrip, err = trans.Prepare("UPDATE trips SET title=?, description=?,hotel=?,stars=?,price=?,date=?,numberOfDays=?,city=?,country=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut edita!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer updateTrip.Close()
	_, err = updateTrip.Exec(title, description, hotel, stars, price, date, days, city, country, param.ByName("tripId"))
	if err != nil {
		fmt.Println(err)
		var message structs.Comment
		message.ID = "yes"
		message.Username = "Nu s-a putut edita!"
		temp.ExecuteTemplate(w, "index.html", message)
		trans.Rollback()
		return
	}
	defer updateTrip.Close()
	var message structs.Comment
	message.ID = "yes"
	message.Username = "Excursie editata cu succes!"
	temp.ExecuteTemplate(w, "index.html", message)
	trans.Commit()

}

func CreateTrip(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		var message structs.Comment
		message.ErrMessage = "Nu ai acces!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	if !session.IsNew && session.Values["Role"].(string) != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu ai acces!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	temp.ExecuteTemplate(w, "createTrip.html", nil)
}

func CreateTripLogic(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	if session.IsNew {
		var message structs.Comment
		message.ErrMessage = "Nu ai acces!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}

	if !session.IsNew && session.Values["Role"].(string) != "AGENCY" {
		var message structs.Comment
		message.ID = "yes"
		message.ErrMessage = "Nu ai acces!"
		temp.ExecuteTemplate(w, "index.html", message)
		return
	}
	var title string = r.FormValue("title")
	var description string = r.FormValue("description")
	var hotel string = r.FormValue("hotel")
	var stars string = r.FormValue("stars")
	var price string = r.FormValue("price")
	if priceFloat, err := strconv.ParseFloat(price, 64); err == nil {
		if priceFloat < 0 {
			temp.ExecuteTemplate(w, "createTrip.html", "Pretul nu poate fi negativ")
			return
		}
	}
	var date string = r.FormValue("date")
	fmt.Println(date)
	var days string = r.FormValue("days")
	var city string = r.FormValue("city")
	var country string = r.FormValue("country")
	if len(title) <= 0 || len(description) <= 0 || len(hotel) <= 0 || len(stars) <= 0 || len(price) <= 0 || len(date) <= 0 || len(days) <= 0 || len(city) <= 0 || len(country) <= 0 {
		temp.ExecuteTemplate(w, "createTrip.html", "Exista campuri necompletate")
		return
	}
	//first image
	file1, photoChar1, err1 := r.FormFile("img1")
	if err1 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", "Ceva nu este in regula cu poza adaugata!")
		return
	}
	defer file1.Close()
	fmt.Println(photoChar1.Filename)
	//add in the directory the file
	photo1, err1 := ioutil.TempFile("assets\\images", "trip-*.png")
	// take the last dash position
	fmt.Println(photo1.Name())
	if err1 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", err1.Error())
		return
	}
	defer photo1.Close()

	photoBytes1, err1 := ioutil.ReadAll(file1)
	if err1 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", err1.Error())
		return
	}
	// write this byte array to our folder
	photo1.Write(photoBytes1)

	//second image
	file2, photoChar2, err2 := r.FormFile("img2")
	if err2 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", "Ceva nu este in regula cu poza adaugata!")
		return
	}
	defer file2.Close()
	fmt.Println(photoChar2.Filename)
	//add in the directory the file
	photo2, err2 := ioutil.TempFile("assets\\images", "trip-*.png")
	// take the last dash position
	fmt.Println(photo2.Name())
	if err2 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", err2.Error())
		return
	}
	defer photo2.Close()

	photoBytes2, err2 := ioutil.ReadAll(file2)
	if err2 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", err2.Error())
		return
	}
	// write this byte array to our folder
	photo2.Write(photoBytes2)

	//third image
	file3, photoChar3, err3 := r.FormFile("img3")
	if err3 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", "Ceva nu este in regula cu poza adaugata!")
		return
	}
	defer file3.Close()
	fmt.Println(photoChar3.Filename)
	//add in the directory the file
	photo3, err3 := ioutil.TempFile("assets\\images", "trip-*.png")
	// take the last dash position
	fmt.Println(photo3.Name())
	if err3 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", err3.Error())
		return
	}
	defer photo3.Close()

	photoBytes3, err3 := ioutil.ReadAll(file3)
	if err3 != nil {
		temp.ExecuteTemplate(w, "createTrip.html", err3.Error())
		return
	}
	// write this byte array to our folder
	photo3.Write(photoBytes3)

	// start the transaction, because all the validations passed at this point
	var trans *sql.Tx
	trans, err := database.Db.Begin()
	if err != nil {
		temp.ExecuteTemplate(w, "createTrip.html", "A aparut o eroare, te rog mai incearca!")
		return
	}
	//this will be ignored in case of a commit
	defer trans.Rollback()

	//insert the trip
	var insertTrip *sql.Stmt
	insertTrip, err = trans.Prepare("INSERT INTO trips (title, description, hotel, stars, price, img1, img2, img3, date, agencyID, numberOfDays,city,country) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);")

	if err != nil {
		temp.ExecuteTemplate(w, "createTrip.html", "Nu s-a putut inregistra2")
		fmt.Println(err)
		trans.Rollback()
		return
	}
	defer insertTrip.Close()

	_, err = insertTrip.Exec(title, description, hotel, stars, price, photo1.Name(), photo2.Name(), photo3.Name(), date, session.Values["Id"], days, city, country)

	if err != nil {
		temp.ExecuteTemplate(w, "registerAgency.html", "Ceva nu a mers cum trebuie!")
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

	m.Username = "Excursia a fost adaugata cu succes"
	m.ID = "yes"
	temp.ExecuteTemplate(w, "index.html", m)

}

func JsonAllTrips(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//select all username of agencies
	sqlQueryA := "SELECT * from trips"
	rows, err := database.Db.Query(sqlQueryA)
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
	trips := make([]*structs.Trip, 0)
	for rows.Next() {
		trip := new(structs.Trip)
		var id, title, description, hotel, stars, price, img1, img2, img3, date, agencyID, numberOfDays, city, country string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&id, &title, &description, &hotel, &stars, &price, &img1, &img2, &img3, &date, &agencyID, &numberOfDays, &city, &country) != nil {
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
		trip.Title = title
		trip.Description = description
		trip.Date = date
		trip.NumberOfDays = numberOfDays
		trip.Hotel = hotel
		trip.Stars = stars
		trip.Price = price
		trip.Path_img1 = img1
		trip.Path_img2 = img2
		trip.Path_img3 = img3
		trip.City = city
		trip.Country = country
		session, _ := Store.Get(r, "session")
		//we send a message if the user is connected so that some buttons will not display
		if !session.IsNew {
			if session.Values["Role"].(string) == "ADMIN" {
				trip.IsTheSame = "yes"
			} else if session.Values["Id"] == agencyID && session.Values["Role"].(string) == "AGENCY" {
				trip.IsTheSame = "yes"
			} else {
				trip.IsTheSame = "no"
			}

			if session.Values["Role"].(string) == "CLIENT" {
				trip.IsClient = "yes"
				trip.ClientID = session.Values["Id"].(string)
			} else {
				trip.IsClient = "no"
				trip.ClientID = "nothing"
			}
		} else {
			session.Options.MaxAge = -1
			session.Save(r, w)
			trip.IsTheSame = "no"
		}
		sqlQuery := "SELECT username from agencies WHERE id=?"
		row := database.Db.QueryRow(sqlQuery, agencyID)
		var username string
		if row.Scan(&username) != nil {
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
		trip.AgencyName = username
		trips = append(trips, trip)

	}
	// make an result of type Trip in order to create a json to send for the html page
	out, err := json.MarshalIndent(trips, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func AllTrips(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	temp.ExecuteTemplate(w, "allTrips.html", nil)
}

func JsonWeather(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	key := os.Getenv("OWM_API_KEY")
	write, err := weather.NewCurrent("C", "RO", key)
	if err != nil {
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
	write.CurrentByName(param.ByName("tripName"))
	out, err := json.MarshalIndent(write, "", "   ")
	if err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func Weather(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	var message structs.Comment
	message.Username = param.ByName("tripName")
	temp.ExecuteTemplate(w, "weather.html", message)
}
