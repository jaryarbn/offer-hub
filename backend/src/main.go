package main

import (
	"fmt"
	"log"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
	"offer-hub/backend/src/ngin"
	rt "offer-hub/backend/src/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

func run() error {
	if err := config.Init(); err != nil {
		return fmt.Errorf("initialize config: %w", err)
	}

	initializedData, err := data.NewData(config.Conf)
	if err != nil {
		return fmt.Errorf("initialize databases: %w", err)
	}
	defer func() {
		if err := initializedData.Close(); err != nil {
			log.Printf("close databases: %v", err)
		}
	}()

	router, err := ngin.CreateGin()
	if err != nil {
		return fmt.Errorf("create Gin engine: %w", err)
	}

	if err := rt.RegisterRouter(router); err != nil {
		return fmt.Errorf("register router: %w", err)
	}

	if err := router.Run(config.Conf.Common.Address()); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}
	return nil
}
