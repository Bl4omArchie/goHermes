package core


type Document struct {
	Title string
	Authors []string
	Download_path string
	Category string
	Document_hash string
	Release_date string
}


type Author struct {
	First_name string
	Last_name string
}
