package db

import (
	_ "fmt"
	_ "os"
	_ "database/sql"
	_ "github.com/lib/pq"
	_ "github.com/Bl4omArchie/ePrint-DB/src/utils"
)


type Papers struct {
	Id int
    Title string
    Link string
    Publication_year string
	Category string
	File_data string
	File_type string
	Page_url string
	Doc_url string
}


// Insert the given paper into the database
func StorePdf(paper Papers) {

}