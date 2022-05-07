package main

import (
	"fmt"
	"log"
	"os"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/infrastructure/weatherservice"
	"weather-alerts-service/internal/logging"
	"weather-alerts-service/internal/usecases"
)

func main() {

	logging.Init()

	opts, err := GetProgramOptions()
	if err != nil {
		log.Fatal("Unable to obtain options for program execution", err)
	}

	mongoClient, err := persistence.NewMongoClient(opts.MongoHost, opts.MongoUser, opts.MongoPassword)
	if err != nil {
		log.Fatal("Unable to create mongo client", err)
	}

	syncAlerts := usecases.SyncAlertDependencies{
		AlertSource:         weatherservice.NewAlertsAdapter("test-application"),
		AlertRepository:     persistence.NewAlertMongoRepository(mongoClient),
		ImportLogRepository: persistence.NewMongoImportLogRepository(mongoClient),
	}

	alertsCount, err := syncAlerts.SyncAlerts()
	if err != nil {
		log.Fatal("Unable to sync alerts.", err)
	}

	fmt.Println("Imported alerts: ", alertsCount)
}

// TODO: will want to retrieve password from a secret file

type ProgramOptions struct {
	MongoHost     string
	MongoUser     string
	MongoPassword string
}

// TODO: Research some go modules for configuration best practices.  However, for now, we will keep this simple

func GetProgramOptions() (ProgramOptions, error) {

	password := os.Getenv("MONGO_PASSWORD")

	/*
		password, err := GetSecret("MONGO_PASSWORDFILE")
		if err != nil {
			return ProgramOptions{}, err
		}*/

	options := ProgramOptions{
		MongoHost:     os.Getenv("MONGO_HOST"),
		MongoUser:     os.Getenv("MONGO_USER"),
		MongoPassword: password,
	}

	return options, nil
}
