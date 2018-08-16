package main

import (
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

func getDB() *pop.Connection {
	db, err := pop.Connect(state.Environment)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Established DB connection.")
	pop.Debug = state.Environment == "development"
	return db
}

func migrateDB() {
	fileMigrator, err := pop.NewFileMigrator("./migrations", state.DB)

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
