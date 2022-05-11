package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"myapp/database"
	"myapp/structs"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func JsonAgency(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	fmt.Println(param.ByName("nameOfAgency"))
	//select in another query the particularities from agencies
	sqlQueryA := "SELECT username,description,email,path_profile_image from agencies where username=?"
	rowA := database.Db.QueryRow(sqlQueryA, param.ByName("nameOfAgency"))
	var username, description, email, path_profile_image string
	//if don't exist => no agency
	if rowA.Scan(&username, &description, &email, &path_profile_image) != nil {
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
	// make an result of type Agency in order to create a json to send for the html page
	fmt.Println(path_profile_image)
	var result structs.Agency
	result.Username = username
	result.Description = description
	result.Email = email
	result.Profile_image = path_profile_image
	// indent de json
	out, err := json.MarshalIndent(result, "", "   ")
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
			message.ErrMessage = "Ceva nu a mers cum trebuie!"
			session.Save(r, w)
			temp.ExecuteTemplate(w, "index.html", nil)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func AgencyPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	var m structs.Comment
	m.Username = param.ByName("nameOfAgency")
	session, _ := Store.Get(r, "session")
	if !session.IsNew {
		m.ID = "yes"
		if m.Username == session.Values["Username"].(string) {
			m.IsTheSame = "yes"
			m.ID = "yes"
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
			message.ErrMessage = "Ceva nu a mers cum trebuie!"
			temp.ExecuteTemplate(w, "index.html", message)
		}
		return
	}
	temp.ExecuteTemplate(w, "agencyPersonalPage.html", m)
}

func JsonAllAgencies(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	//select all username of agencies
	sqlQueryA := "SELECT id,username,path_profile_image from agencies"
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
	agencies := make([]*structs.Agency, 0)
	for rows.Next() {
		agency := new(structs.Agency)
		var id, username, path_profile_image string
		//if the username and path to profile is not the same redirect
		if rows.Scan(&id, &username, &path_profile_image) != nil {
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
		agency.ID = id
		agency.Profile_image = path_profile_image
		agency.Username = username
		session, _ := Store.Get(r, "session")
		//we send a message if the user is connected so that some buttons will not display
		if !session.IsNew {
			if session.Values["Role"] == "ADMIN" {
				agency.Is_admin = "yes"
			}
		} else {
			session.Options.MaxAge = -1
			session.Save(r, w)
		}
		agencies = append(agencies, agency)
	}
	// make an result of type Agency in order to create a json to send for the html page
	out, err := json.MarshalIndent(agencies, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func AllAgencies(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	temp, _ = template.ParseGlob("templates/*.html")
	session, _ := Store.Get(r, "session")
	var m structs.Comment
	if !session.IsNew {
		m.ID = "yes"
	} else {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	temp.ExecuteTemplate(w, "allAgencies.html", m)
}
