package core

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Source interface {
	Init(engine *Engine) error
	Fetch(engine *Engine) error
}

type Engine struct {
	Log *Log
	SqliteDb *gorm.DB
	DatabaseName string
	NumWorkersPools int
	Sources []Source
}

func StartEngine() (error) {
	engine, err := CreateEngineInstance()
	if (err != nil) {
		CreateLogReport("Failed to ceate engine instance", engine.Log)
		return err
	}
	defer ExitEngineInstance(engine)

	fmt.Println(engine.SqliteDb.Statement.Vars...)

	if err := MigrateSqliteDatabase(engine, &Document{}); err != nil {
		CreateLogReport("Failed to migrate database", engine.Log)
		return err
	}

	for _, src := range engine.Sources {
		src.Init(engine)
		src.Fetch(engine)
	}

	return nil
}

func CreateEngineInstance() (*Engine, error) {
	// Temporary setup of engine parameters
	databaseName := "core/eprint.db"
	numWorkersPools := 100
	
	engine := &Engine{DatabaseName: databaseName, NumWorkersPools: numWorkersPools}

	CreateLogChannel(engine)
	go ListenerLogFile(engine.Log)

	if _, err := os.Stat(databaseName); os.IsNotExist(err) {
		file, err := os.Create(databaseName)
		if err != nil {
			CreateLogReport("Failed at attempting to create the database", engine.Log)
			return nil, err
		}
		defer file.Close()
	}

	err := OpenSqliteDatabase(engine)
	if err != nil {
		CreateLogReport("Failed to open database", engine.Log)
		return nil, err
	}

	engine.Sources = []Source{
		NewEprintSource(),
	}

	return engine, nil
}

func ExitEngineInstance(engine *Engine) (error) {
	close(engine.Log.logChannel)
	err := CloseSqliteDatabase(engine)
	return err
}
