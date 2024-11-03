package db

import (
    "database/sql"
	"fmt"
	"os"
    _ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
	//"github.com/Bl4omArchie/ePrint-DB/src/api"
)

type Paper struct {
	Id int64
    Title string
    Link string
    Publication_date int64
}

type Database struct {
	ConnectionChain string
	SqlDatabase *sql.DB
	Name string
}

func ConnectDatabase() (*Database, error) {
	err := godotenv.Load()
	utils.CheckError(err)

	var (
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
		user = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
	)

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	utils.CheckError(err)

	err = db.Ping()
	utils.CheckError(err)

	fmt.Println("> Connected !")

    return &Database{
        ConnectionChain: psqlconn,
        SqlDatabase: db,
		Name: dbname,
    }, nil
}

func DisconnectDatabase(db *Database) {
	err := db.SqlDatabase.Close()
	utils.CheckError(err)
	fmt.Println("> Disconnected !")
}