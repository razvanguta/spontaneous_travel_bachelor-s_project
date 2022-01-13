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
	database.Db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/database_architecture")
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
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	http.ListenAndServe(":8080", router)

	database.CloseConnectionToDatabase(database.Db, err)
}
