package main

import (
	"context"
	"log"

	"github.com/redhat-openshift-ecosystem/knex/cmd/knex/root"
	"github.com/spf13/viper"
)

func main() {
	appContext := context.Background()
	appConfig := viper.New()

	entrypoint := root.NewCommand(
		appContext,
		appConfig,
	)

	if err := entrypoint.Execute(); err != nil {
		log.Fatal(err)
	}
}
