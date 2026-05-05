package models

import "gorm.io/gorm"


type Author struct {
	gorm.Model
	FirstName	string     `gorm:"not null"`
	LastName 	string     `gorm:"not null"`
	Documents	[]Document `gorm:"many2many:Author_authors;"`
}


func NewAuthor(first, last string, docs []Document) *Author {
	return &Author{
		FirstName: first,
		LastName: last,
		Documents: docs,
	}
}
