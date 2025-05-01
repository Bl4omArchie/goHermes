package data

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	ID int				`gorm:"primaryKey"`
	Title string		`gorm:"unique;not null"`
	Authors []Author	`gorm:"many2many:author_documents;not null"`
	Filepath string		`gorm:"unique;not null"`
	Hash string			`gorm:"unique;not null"`
	Release string		`gorm:"not null"`
	License string
	Source string		`gorm:"not null"`
}

type Author struct {
	gorm.Model
	ID int 					`gorm:"primaryKey"`
	FirstName string		`gorm:"not null"`
	LastName string
	Documents []Document	`gorm:"many2many:author_documents"`
}


func CreateDocument(title string, authors []Author, filepath string, hash string, release string, license string, source string) (*Document) {
	return &Document{
		Title:    title,
		Authors:  authors,
		Filepath: filepath,
		Hash:     hash,
		Release:  release,
		License:  license,
		Source:    source,
	}
}

func CreateAuthor(firstName string, lastName string) (*Author) {
	return &Author {
		FirstName: firstName,
		LastName: lastName,
	}
}
