package db

import (
	_ "fmt"
	_ "os"
	_ "database/sql"
	_ "github.com/lib/pq"
	_ "github.com/Bl4omArchie/ePrint-DB/src/utils"
)


type Papers struct {
    Title string
    Link string
    Publication_year int
	Category string
	File_data string
}



// Insert the given paper into the database
func StorePdf(paper Papers) {

}