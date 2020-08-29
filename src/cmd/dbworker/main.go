package main

import (
	"bdim/src/internal/dbworker"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := dbworker.NewConfig()
	err := config.ParseFlags()
	if err != nil {
		dbworker.Logger.Fatal(err)
	}
	dbworker.Logger.Printf("Version: %s", dbworker.Ver)

	consumer, err := dbworker.NewConsumer(config)
	if err != nil {
		dbworker.Logger.Fatal(err)
	}
	loader, err := dbworker.NewLoader(config)
	if err != nil {
		dbworker.Logger.Fatal(err)
	}

	consumer.Batches(func(events [][]byte) error {
		data, err := dbworker.ParseJson(events)
		if err != nil {
			return err
		}
		rowsAffected, err := loader.Upsert(data)
		if err != nil {
			return err
		}
		dbworker.Logger.Printf("Affected rows: %d", rowsAffected)
		return nil
	})

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<- wait

	err = consumer.Close()
	if err != nil {
		dbworker.Logger.Printf("Failed to close consumer: %v", err)
	}
}