package main

import (
	"log"

	"github.com/Xacor/gophermart/internal/app"
	"github.com/Xacor/gophermart/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
