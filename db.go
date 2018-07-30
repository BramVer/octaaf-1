package main

import (
	"log"

	"github.com/gobuffalo/pop"
)

// DB Global Connection
var DB *pop.Connection

func connectDB() {
	var err error
	DB, err = pop.Connect(OctaafEnv)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = OctaafEnv == "development"
}

func migrateDB() {
	fileMigrator, err := pop.NewFileMigrator("./migrations", DB)

	if err != nil {
		log.Panic(err)
	}

	err = fileMigrator.Up()

	if err != nil {
		log.Panic(err)
	}
}
