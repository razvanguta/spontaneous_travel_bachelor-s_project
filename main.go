package main

import (
	"database/sql"
	"myapp/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_architecture")
	database.ConnectToDatabase(db, err)

	database.CloseConnectionToDatabase(db, err)
}
