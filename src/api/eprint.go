package api


import (
	_ "github.com/cavaliergopher/grab/v3"
	"github.com/Bl4omArchie/ePrint-DB/src/db"
)

var (
	url = "https://eprint.iacr.org/"
	
	categories = []string {
		"Applications",
		"Cryptographic protocols",
		"Foundations",
		"Implementation",
		"Secret-key cryptography",
		"Public-key cryptography",
		"Attacks and cryptanalysis",
	}

	years = []int {2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015, 2014, 2013, 2012, 2010, 2009, 2008, 2007, 2006, 2005, 2004, 2003, 2002, 2001, 2000, 1999, 1998, 1997, 1996}
)

func DownloadPapersByCategories(year int) {
		
}

func DownloadPapersByYears(year int) {
		
}

func DownloadPapers(years []int) {
		
}

func StartApplication() {
	database := db.ConnectDatabase()

	defer db.DisconnectDatabase(database)
}