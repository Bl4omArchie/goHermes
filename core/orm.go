package core


import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Document struct {
	gorm.Model
	ID int				`gorm:"primaryKey"`
	Title string		`gorm:"unique;not null"`
	Authors []Author	`gorm:"many2many:author_documents;not null"`
	Filepath string		`gorm:"unique;not null"`
	Url string			`gorm:"not null"`
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


func OpenSqliteDatabase(errChannel *ErrorChannel) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		CreateErrorReport("Failed to connect to sqlite database", errChannel)
		return nil, err
	}
	return db, err
}

func MigrateSqliteDatabase(db *gorm.DB, errChannel *ErrorChannel, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		CreateErrorReport("Migration failed", errChannel)
		return err
	}
	return nil
}
