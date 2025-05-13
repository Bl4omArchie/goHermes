package core


import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Document struct {
	gorm.Model
	Title string		`gorm:"unique;not null"`
	Authors []Author	`gorm:"many2many:author_documents;not null"`
	Filepath string		`gorm:"unique;not null"`
	Url string			`gorm:"unique;not null"`
	Hash string			`gorm:"unique;not null"`
	Release string		`gorm:"not null"`
	License string
	Source string		`gorm:"not null"`
}

type Author struct {
	gorm.Model
	FirstName string		`gorm:"not null"`
	LastName string
	Documents []Document	`gorm:"many2many:author_documents"`
}


func OpenSqliteDatabase(databaseName string, log *Log) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		CreateLogReport("Failed to connect to sqlite database", log)
		return nil, err
	}
	return db, err
}

func MigrateSqliteDatabase(db *gorm.DB, log *Log, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		CreateLogReport("Migration failed", log)
		return err
	}
	return nil
}

func InsertDocument(doc *Document, authors *[]Author, engineInstance *Engine) error {
	err := engineInstance.SqliteDb.Create(doc).Error
	if err != nil {
		CreateLogReport("Failed to insert document into database", engineInstance.Log)
		return err
	}

	for _, author := range *authors {
		var existingAuthor Author
		if err := engineInstance.SqliteDb.Where("first_name = ? AND last_name = ?", author.FirstName, author.LastName).First(&existingAuthor).Error; err != nil {
			if err := engineInstance.SqliteDb.Create(&author).Error; err != nil {
				CreateLogReport("Failed to insert author into database", engineInstance.Log)
				return err
			}
			doc.Authors = append(doc.Authors, author)
		} else {
			doc.Authors = append(doc.Authors, existingAuthor)
		}
	}
	return nil
}

func InsertAuthor(db *gorm.DB, author *Author, log *Log) error {
	err := db.Create(author).Error
	if err != nil {
		CreateLogReport("Failed to insert author into database", log)
		return err
	}
	return nil
}
