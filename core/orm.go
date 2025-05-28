package core

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
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
    Source   string   `gorm:"not null"`
}

type Author struct {
    gorm.Model
    FirstName string     `gorm:"not null"`
    LastName  string     `gorm:"not null"`
    Documents []Document `gorm:"many2many:document_authors;"`
}

func OpenSqliteDatabase(databaseName string, log *Log) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		CreateLogReport("Failed to connect to sqlite database", log)
		return nil, err
	}
	return db, err
}

func MigrateSqliteDatabase(engine *Engine, tables ...any) error {
	if err := engine.SqliteDb.AutoMigrate(tables...).Error; err != nil {
		CreateLogReport("Migration failed", engine.Log)
		return fmt.Errorf("failed to insert table: %w", err)
	}
	return nil
}

func InsertTable(engine *Engine, tables ...any) error {
	for _, table := range tables {
		if err := engine.SqliteDb.Create(table).Error; err != nil {
			CreateLogReport("Insert failed", engine.Log)
			return fmt.Errorf("failed to insert table: %w", err)
		}
	}
	return nil
}
