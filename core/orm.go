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
	Url string			`gorm:"unique;not null"`
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


func OpenSqliteDatabase(logChannel *LogChannel) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("eprint.db"), &gorm.Config{})
	if err != nil {
		CreateLogReport("Failed to connect to sqlite database", logChannel)
		return nil, err
	}
	return db, err
}

func MigrateSqliteDatabase(db *gorm.DB, logChannel *LogChannel, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		CreateLogReport("Migration failed", logChannel)
		return err
	}
	return nil
}

func InsertDocument(db *gorm.DB, doc *Document, authors *[]Author, logChannel *LogChannel) error {
	err := db.Create(doc).Error
	if err != nil {
		CreateLogReport("Failed to insert document into database", logChannel)
		return err
	}

	for _, author := range *authors {
		var existingAuthor Author
		if err := db.Where("first_name = ? AND last_name = ?", author.FirstName, author.LastName).First(&existingAuthor).Error; err != nil {
			if err := db.Create(&author).Error; err != nil {
				CreateLogReport("Failed to insert author into database", logChannel)
				return err
			}
			doc.Authors = append(doc.Authors, author)
		} else {
			doc.Authors = append(doc.Authors, existingAuthor)
		}
	}
	return nil
}

func InsertAuthor(db *gorm.DB, author *Author, logChannel *LogChannel) error {
	err := db.Create(author).Error
	if err != nil {
		CreateLogReport("Failed to insert author into database", logChannel)
		return err
	}
	return nil
}
