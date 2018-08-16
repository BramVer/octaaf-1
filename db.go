package main

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

var DB *pop.Connection

func initDB() {
	// Don't refer to state.Environment as this function is called before it's available
	var err error
	DB, err = pop.Connect(envy.Get("GO_ENV", "development"))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Established DB connection.")
	pop.Debug = state.Environment == "development"
}

func migrateDB() {
	fileMigrator, err := pop.NewFileMigrator("./migrations", DB)

	if err != nil {
		log.Panic(err)
	}

	fileMigrator.Status()

	err = fileMigrator.Up()

	if err != nil {
		log.Panic(err)
	}
	log.Info("Finished DB migrations")
}
