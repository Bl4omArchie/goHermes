package hermes

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Title    string   `gorm:"unique;not null"`
	Authors  []Author `gorm:"many2many:document_authors;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Filepath string   `gorm:"unique;not null"`
	Filetype string   `gorm:"not null"`
	Url      string   `gorm:"unique;not null"`
	Hash     string   `gorm:"not null"`
	Release  string   `gorm:"not null"`
	License  string
	Source   string `gorm:"not null"`
}

type Author struct {
	gorm.Model
	FirstName string     `gorm:"not null"`
	LastName  string     `gorm:"not null"`
	Documents []Document `gorm:"many2many:document_authors;"`
}

func OpenSqliteDatabase(engine *Engine) error {
	db, err := gorm.Open(sqlite.Open(engine.DatabaseName), &gorm.Config{})
	if err != nil {
		CreateLogReport("Failed to connect to sqlite database", engine.Log)
		return err
	} else {
		engine.SqliteDb = db
	}
	return nil
}

func CloseSqliteDatabase(engine *Engine) error {
	sqlDB, err := engine.SqliteDb.DB()
	if err != nil {
		CreateLogReport("Failed to get sql.DB instance from gorm.DB", engine.Log)
		return err
	}

	if err := sqlDB.Close(); err != nil {
		CreateLogReport("Failed to close sqlite database", engine.Log)
		return err
	}

	return nil
}

func MigrateSqliteDatabase(engine *Engine, tables ...any) error {
	if err := engine.SqliteDb.AutoMigrate(tables...); err != nil {
		CreateLogReport("Migration failed", engine.Log)
		return err
	}
	return nil
}

func InsertTable(engine *Engine, tables ...any) error {
	for _, table := range tables {
		if err := engine.SqliteDb.Create(table).Error; err != nil {
			CreateLogReport(fmt.Sprintf("Insert table failed : %v", err), engine.Log)
			return err
		}
	}
	return nil
}
