package main

import (
	"database/sql"
	"myapp/database"
	"myapp/handlers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	database.Db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/database_architecture")
	database.ConnectToDatabase(database.Db, err)

	//http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/register", handlers.RegisterClientPage)
	http.HandleFunc("/registerClientLogic", handlers.RegisterClientLogic)
	http.HandleFunc("/emailVerification", handlers.EmailVerification)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/loginLogic", handlers.LoginLogic)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/deleteMyself", handlers.DeleteMyself)
	http.HandleFunc("/", handlers.HomePage)
	http.ListenAndServe(":8080", nil)

	database.CloseConnectionToDatabase(database.Db, err)
}
