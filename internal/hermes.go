package hermes

import (
	"fmt"
	"sync"

	"github.com/Bl4omArchie/goHermes/models"
)

type HermesApp struct {
	Logging  *SlogWrapper
	Database *HermesDatabase
	Workers  int
	Sources  []string
}

func NewHermesApp(log *SlogWrapper, db *HermesDatabase, workers int, sources []string) *HermesApp {
	return &HermesApp{
		Logging:  log,
		Database: db,
		Workers:  workers,
		Sources:  sources,
	}
}

func CreateEngine(sources []string, config *HermesConfig) (*HermesApp, error) {
	log := InitializeSlogWrapper(config.LoggingDir, config.LoggingFormat)

	db, err := CreateHermesDatabase(config.DatabaseName)
	if err != nil {
		return nil, err
	}

	return NewHermesApp(log, db, config.Workers, sources), nil
}

func (ha *HermesApp) Start() error {
	defer ha.Exit()

	if err := ha.Database.Db.AutoMigrate(&models.Document{}); err != nil {
		return err
	}

	// var wg sync.WaitGroup
	// for _, src := range ha.Sources {
	// 	wg.Add(1)
	// 	go func() {
	// 		src.Fetch(engine)
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()
	// fmt.Println("Every document has been downloaded")

	return nil
}

func (ha *HermesApp) Exit() error {
	err := ha.Database.Db.Close(ha)
	return err
}
