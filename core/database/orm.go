package database


import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/Bl4omArchie/eprint-DB/core/utility"
)


func OpenSqliteDatabase(errChannel *utility.ErrorChannel) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		utility.CreateErrorReport("Failed to connect to sqlite database", errChannel)
		return nil, err
	}
	return db, err
}

func MigrateSqliteDatabase(db *gorm.DB, errChannel *utility.ErrorChannel, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		utility.CreateErrorReport("Migration failed", errChannel)
		return err
	}
	return nil
}

func CreateSqliteDatabase(db *gorm.DB, errChannel *utility.ErrorChannel, models ...any) error {
	err := db.Create(models...)
	if err != nil {
		utility.CreateErrorReport("Base creation failed", errChannel)
		return err
	}
	return nil
}