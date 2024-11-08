package db

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)


type Database struct {
	ConnectionChain string
	SqlDatabase *sql.DB
	Name string
}


func ConnectDatabase() (*Database) {
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
	utils.CheckErrorQuit(err)

	err = db.Ping()
	utils.CheckErrorQuit(err)

	fmt.Println("\033[32m> Connected !\033[0m\n")

    return &Database{
        ConnectionChain: psqlconn,
        SqlDatabase: db,
		Name: dbname,
    }
}

func DisconnectDatabase(db *Database) {
	err := db.SqlDatabase.Close()
	utils.CheckError(err)
	fmt.Println("\033[32m> Disconnected !\033[0m")
}