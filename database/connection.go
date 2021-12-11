package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func ConnectToDatabase(db *sql.DB, err error) {

	//check the arguments are good
	if err != nil {
		fmt.Println("The arguments in sql.Open are not good")
		panic(err.Error()) //stop the goroutine
	}

	err = db.Ping()
	if err == nil {
		fmt.Println("Connection to database is open!")
	} else {
		fmt.Println("No ping from database!")
		panic(err.Error())
	}
}

func CloseConnectionToDatabase(db *sql.DB, err error) {
	fmt.Println("Connection to database is close!")
	db.Close()
}
