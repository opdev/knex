package main

import (
	"context"
	"log"

	"github.com/opdev/knex/cmd/knex/root"
)

func main() {
	entrypoint := root.NewCommand(
		context.Background(),
	)

	if err := entrypoint.Execute(); err != nil {
		log.Fatal(err)
	}
}
