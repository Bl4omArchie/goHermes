package main

import (
	"fmt"

	"github.com/Bl4omArchie/goHermes/internal"

)

func main() {
	fmt.Println("Launching downloading engine...")
	engine, _ := hermes.CreateEngineInstance("hermes.db", 50, hermes.NewFreeHavenSource(), hermes.NewEprintSource())
	hermes.StartEngine(engine)
}
