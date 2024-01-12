package main

import (
	"fmt"
	"log"

	"github.com/erupshis/key_keeper/internal/agent/controller"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
)

func main() {
	fmt.Println("Простое CLI-приложение")

	inMemoryStorage := inmemory.NewStorage()
	mainController := controller.NewController(inMemoryStorage)

	if err := mainController.Serve(); err != nil {
		log.Fatalf("problem with controller: %v", err)
	}

}
