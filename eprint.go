package main

import (
	"github.com/Bl4omArchie/ePrint-DB/src/db"
)



func main() {
	database, err = db.ConnectDatabase()

	defer db.DisconnectDatabase(database)
}