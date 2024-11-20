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


func ConnectDatabase(ac *utils.AlertChannel) (*Database) {
	err := godotenv.Load()
	utils.CheckAlertError(err, 0xc6, "Incorrect credential for database.", ac)

	var (
		host = os.Getenv("HOST")
		port = os.Getenv("PORT")
		user = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname = os.Getenv("NAME")
	)

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	utils.CheckErrorQuit(err)

	err = db.Ping()
	utils.CheckErrorQuit(err)

	fmt.Println("\033[32m> Connected !\033[0m")

    return &Database{
        ConnectionChain: psqlconn,
        SqlDatabase: db,
		Name: dbname,
    }
}

func DisconnectDatabase(ac *utils.AlertChannel, db *Database) {
	err := db.SqlDatabase.Close()
	utils.CheckAlertError(err, 0xc6, "Error while deconnecting.", ac)
	fmt.Println("\033[32m> Disconnected !\033[0m")
}