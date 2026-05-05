package hermes

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)



type HermesDatabase struct {
	Path 	string
	Db		*gorm.DB
}


func NewHermesDatabase(path string, db *gorm.DB) *HermesDatabase {
	return &HermesDatabase{
		Path: path,
		Db: db,
	}
}


func CreateHermesDatabase(name string) (*HermesDatabase, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}
	
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return NewHermesDatabase(name, db), nil
}


func (hd *HermesDatabase) Close() error {
	sqlDB, err := hd.Db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}


func (hd *HermesDatabase) InsertTable(tables ...any) error {
	for _, table := range tables {
		if err := hd.Db.Create(table).Error; err != nil {
			return err
		}
	}
	return nil
}
