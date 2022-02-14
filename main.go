package main

import (
	"database/sql"
	"myapp/database"
	"myapp/handlers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var err error
	database.Db, err = sql.Open("mysql", "root:mypass@tcp(localhost:3306)/database_architecture")
	database.ConnectToDatabase(database.Db, err)
	router := httprouter.New()
	router.GET("/", handlers.HomePage)
	router.GET("/register", handlers.RegisterClientPage)
	router.POST("/registerClientLogic", handlers.RegisterClientLogic)
	router.GET("/registerAgency", handlers.RegisterAgency)
	router.POST("/registerAgencyLogic", handlers.RegisterAgencyLogic)
	router.GET("/agencyPage/:nameOfAgency", handlers.AgencyPage)
	router.GET("/jsonAgency/:nameOfAgency", handlers.JsonAgency)
	router.GET("/emailVerificationPage", handlers.EmailVerificationPage)
	router.POST("/emailVerification", handlers.EmailVerification)
	router.GET("/login", handlers.Login)
	router.POST("/loginLogic", handlers.LoginLogic)
	router.GET("/logout", handlers.Logout)
	router.GET("/deleteMyself", handlers.DeleteMyself)
	router.GET("/editDescriptionAgency", handlers.EditDescriptionAgency)
	router.POST("/editDescriptionAgencyPut", handlers.EditDescriptionAgencyPut)
	router.GET("/myPersonalPage", handlers.PersonalPage)
	router.GET("/jsonAllAgencies", handlers.JsonAllAgencies)
	router.GET("/allAgencies", handlers.AllAgencies)
	router.GET("/changePassword", handlers.ChangePassword)
	router.POST("/changePasswordLogic", handlers.ChangePasswordLogic)
	router.GET("/passwordReset", handlers.PasswordReset)
	router.POST("/passwordResetLogic", handlers.PasswordResetLogic)
	router.GET("/createTrip", handlers.CreateTrip)
	router.POST("/createTripLogic", handlers.CreateTripLogic)
	router.GET("/allTrips", handlers.AllTrips)
	router.GET("/jsonAllTrips", handlers.JsonAllTrips)
	router.GET("/deleteTrip/:tripId", handlers.DeleteTrip)
	router.GET("/jsonUpdateTrip/:tripId", handlers.JsonUpdateTrip)
	router.GET("/updateTripPage/:tripId", handlers.UpdateTripPage)
	router.POST("/updateTripLogic/:tripId", handlers.UpdateTripLogic)
	router.GET("/jsonWeather/:tripName", handlers.JsonWeather)
	router.GET("/weather/:tripName", handlers.Weather)
	router.GET("/seeReviews/:nameOfAgency", handlers.SeeReviews)
	router.POST("/putReview/:userId/:agencyId", handlers.PutReview)
	router.GET("/jsonReview/:agencyId", handlers.ReviewJson)
	router.GET("/deleteReview/:clientId/:agencyId/:date", handlers.DeleteReview)
	router.GET("/deleteAgency/:agencyId", handlers.DeleteAgency)
	router.GET("/deleteClient", handlers.DeleteClient)
	router.GET("/deleteClientPage", handlers.DeleteAgencyPage)
	router.POST("/addCart/:tripId/:clientId", handlers.AddToCart)
	router.GET("/jsonSeeCart", handlers.JsonSeeCart)
	router.GET("/seeCart", handlers.SeeCart)
	router.GET("/outFromCart/:cartId", handlers.OutFromCart)
	router.GET("/addMoney", handlers.AddMoney)
	router.POST("/addMoneyLogic", handlers.AddMoneyLogic)
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	http.ListenAndServe(":8080", router)

	database.CloseConnectionToDatabase(database.Db, err)
}
