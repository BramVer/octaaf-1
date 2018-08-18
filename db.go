package main

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

var DB *pop.Connection

func initDB() {
	// pop requires a database.yml file
	// This yaml file refers to the DATABASE_URL environment variable as the uri
	// So we set this env variable, so that database.yml can point to it.
	// I hope this changes in the future
	envy.Set("DATABASE_URL", settings.Database.Uri)
	var err error
	DB, err = pop.Connect(settings.Environment)

	if err != nil {
		log.Fatal("Couldn't establish a database connection: ", err)
	}

	log.Info("Established DB connection.")
	pop.Debug = settings.Environment == "development"
}

func migrateDB() {
	fileMigrator, err := pop.NewFileMigrator("./migrations", DB)

	if err != nil {
		log.Fatal(err)
	}

	fileMigrator.Status()

	err = fileMigrator.Up()

	if err != nil {
		log.Fatal(err)
	}
	log.Info("Finished DB migrations")
}
