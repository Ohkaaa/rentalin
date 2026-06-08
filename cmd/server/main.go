package main

import (
	"log"
	"rentalin/config"
	bootstrapp "rentalin/internal/bootstrap"
)

func main() {
	cfg := config.LoadConfig()

	e, db := bootstrapp.NewApp(cfg)

	defer db.Close()

	log.Println("server running on port", cfg.ServerPort)

	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
