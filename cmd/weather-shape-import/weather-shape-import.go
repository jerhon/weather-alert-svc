package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/internal/infrastructure/nwsshapefiles"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/logging"
	"weather-alerts-service/internal/usecases"
)

// TODO: will want to retrieve password from a secret file
// TODO: Research some go modules for configuration best practices.  However, for now, we will keep this simple

func getShapeFile(relativePath string) (string, string, error) {

	fullPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", "", err
	}
	_, fileName := filepath.Split(fullPath)
	shapeFileName := fileName[:(len(fileName)-4)] + ".shp"

	return fullPath, shapeFileName, nil
}

func main() {

	logging.Init()

	countiesFile := flag.String("counties", "", "Zip file containing the county shapes.")

	flag.Parse()

	if *countiesFile == "" {
		flag.Usage()
		return
	}

	opts, _ := GetProgramOptions()

	fullPath, shapeFile, err := getShapeFile(*countiesFile)

	countyShapeMapper := new(nwsshapefiles.CountyShapeAdapter)

	shapeReader := nwsshapefiles.DomainShapeReader[domain.County]{
		ZipFilePath:   fullPath,
		ShapeFileName: shapeFile,
		ShapeAdapter:  countyShapeMapper,
	}

	mongoClient, err := persistence.NewMongoClient(opts.MongoHost, opts.MongoUser, opts.MongoPassword)
	if err != nil {
		log.Fatal("Unable to create mongo client", err)
	}

	repository := persistence.MongoCountiesRepository{
		Client: mongoClient,
	}

	syncCounties := usecases.ImportCountyDeps{
		Repository:      &repository,
		ShapefileReader: &shapeReader,
	}

	err = syncCounties.SyncCounties()
	if err != nil {
		log.Fatal("Unable to sync counties in MongoDB.", err)
	}
}

type ProgramOptions struct {
	MongoHost     string
	MongoUser     string
	MongoPassword string
}

func GetProgramOptions() (ProgramOptions, error) {

	password := os.Getenv("MONGO_PASSWORD")

	options := ProgramOptions{
		MongoHost:     os.Getenv("MONGO_HOST"),
		MongoUser:     os.Getenv("MONGO_USER"),
		MongoPassword: password,
	}

	return options, nil
}
