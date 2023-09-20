package main

import (
	"context"
	"log"

	"github.com/opdev/knex/cmd/knex/root"
)

func main() {
	entrypoint := root.NewCommand()

	ctx := context.Background()
	if err := entrypoint.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
