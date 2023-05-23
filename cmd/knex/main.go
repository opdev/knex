package main

import (
	"context"
	"log"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/redhat-openshift-ecosystem/knex/cmd/knex/root"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	l.SetLevel(logrus.DebugLevel)

	appConfig := viper.New()

	logger := logrusr.New(l)
	appContext := logr.NewContext(context.Background(), logger)

	entrypoint := root.NewCommand(
		appContext,
		appConfig,
	)

	if err := entrypoint.Execute(); err != nil {
		log.Fatal(err)
	}
}
