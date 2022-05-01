package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/infrastructure/shapefiles"
	"weather-alerts-service/internal/logging"
	"weather-alerts-service/internal/usecases"
	"weather-alerts-service/pkg/sliceutils"
)

// TODO: will want to retrieve password from a secret file
// TODO: Research some go modules for configuration best practices.  However, for now, we will keep this simple
// TODO: Maybe use IoC for some of these dependencies?

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

	if len(os.Args) < 2 {
		fmt.Errorf("Usage " + os.Args[0] + " [county|marine] file1 file2...fileN")
	}

	command := os.Args[1]
	files := os.Args[2:]

	opts, _ := GetProgramOptions()

	mongoClient, err := persistence.NewMongoClient(opts.MongoHost, opts.MongoUser, opts.MongoPassword)
	if err != nil {
		log.Fatal("Unable to create mongo client", err)
	}

	readers := sliceutils.MapFunc[string, shapefiles.DomainShapeReader](files, MapFilePathToShapeFileReader)

	sync := usecases.SyncShapeData{
		ShapeType:        command,
		MongoClient:      mongoClient,
		ShapefileReaders: readers,
	}

	err = sync.SyncDomain()
	if err != nil {
		log.Fatal("Unable to sync data.", err)
	}
}

func MapFilePathToShapeFileReader(filePath string) shapefiles.DomainShapeReader {

	fullPath, shapeFile, err := getShapeFile(filePath)
	if err != nil {
		log.Fatal("Unable to create shape reader from file ", err)
	}

	return shapefiles.DomainShapeReader{
		ZipFilePath:   fullPath,
		ShapeFileName: shapeFile,
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
