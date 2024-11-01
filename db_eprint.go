package main

import (
    "database/sql"
	"fmt"
	"os"
    _ "github.com/lib/pq"
	"github.com/joho/godotenv"
)


func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
	err := godotenv.Load()
	CheckError(err)

	var (
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
		user = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
	)

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	fmt.Println("> Connected")
}