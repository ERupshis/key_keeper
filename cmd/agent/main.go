package main

import (
	"log"

	"github.com/erupshis/key_keeper/internal/agent/controller"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
)

func main() {
	inMemoryStorage := inmemory.NewStorage()
	mainController := controller.NewController(inMemoryStorage)

	if err := mainController.Serve(); err != nil {
		log.Fatalf("problem with controller: %v", err)
	}

}
