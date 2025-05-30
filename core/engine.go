package core

import (
	"gorm.io/gorm"
)


type Engine struct {
	Log *Log
	SqliteDb *gorm.DB
	DatabaseName string
	NumWorkersPools int
}

func StartEngine() {
	engineInstance, err := CreateEngineInstance()
	if (err != nil) {
		return
	}

	MigrateSqliteDatabase(engineInstance, &Document{}, &Author{})

	eprint := InitEprint(engineInstance)
	DownloadEprint(eprint, engineInstance)

	ExitEngineInstance(engineInstance)
}

func CreateEngineInstance() (*Engine, error) {
	databaseName := "core/eprint.db"
	numWorkersPools := 100

	log := CreateLogChannel()
	go ListenerLogFile(log)

	database, err := OpenSqliteDatabase(databaseName, log)
	if (err != nil) {
		CreateLogReport("Can't open database", log)
		return nil, err
	}

	return &Engine {
		Log: log,
		SqliteDb: database,
		DatabaseName: databaseName,
		NumWorkersPools: numWorkersPools,
	}, nil
}

func ExitEngineInstance(engineInstance *Engine) {
	// TODO : close DB
	close(engineInstance.Log.logChannel)
}
