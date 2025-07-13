package core

import (
	"os"
	"fmt"
	"sync"

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

func CreateEngineInstance(databaseName string, workers int, sources... Source) (*Engine, error) {	
	engine := &Engine{DatabaseName: databaseName, NumWorkersPools: workers}

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

	engine.Sources = make([]Source, 0, len(sources))
	for _, source := range sources {
		if err := source.Init(engine); err != nil {
			CreateLogReport(fmt.Sprintf("Failed to initialize source: %v", err), engine.Log)
			return nil, err
		}
		engine.Sources = append(engine.Sources, source)
	}

	return engine, nil
}

func StartEngine(engine *Engine) (error) {
	defer ExitEngineInstance(engine)

	if err := MigrateSqliteDatabase(engine, &Document{}); err != nil {
		CreateLogReport("Failed to migrate database", engine.Log)
		return err
	}

	var wg sync.WaitGroup
	for _, src := range engine.Sources {
		wg.Add(1)
		go func() {
			src.Fetch(engine)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Every document has been downloaded")

	return nil
}

func ExitEngineInstance(engine *Engine) (error) {
	close(engine.Log.logChannel)
	err := CloseSqliteDatabase(engine)
	return err
}
