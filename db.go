package main

import (
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

// DB Global Connection
var DB *pop.Connection

func initDB() {
	var err error
	DB, err = pop.Connect(OctaafEnv)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Established DB connection.")
	pop.Debug = OctaafEnv == "development"

	migrateDB()
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
