package usecases

import (
	"go.mongodb.org/mongo-driver/mongo"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/infrastructure/shapefiles"
	"weather-alerts-service/internal/logging"
	"weather-alerts-service/pkg/sliceutils"
)

type SyncShapeData struct {
	ShapefileReaders []shapefiles.DomainShapeReader
	ShapeType        string
	MongoClient      *mongo.Client
}

func (deps *SyncShapeData) getShapefileEntries() ([]domain.DomainShape, error) {

	allShapes := make([]domain.DomainShape, 0)

	for _, reader := range deps.ShapefileReaders {
		shapes, err := reader.GetAllShapes()
		if err != nil {
			return nil, err
		}

		for _, shape := range shapes {
			allShapes = append(allShapes, shape)
		}
	}
	return allShapes, nil
}

func (deps *SyncShapeData) syncDomainTypesInDatabase(shapes []domain.DomainShape) error {

	// TODO: do this with a non-destructive with an upsert so counties can be updated in real time. However, this is likely a one time import as the boundaries do not change often enough
	// Another option would be to create them in a new collection, and then after that's done swap the collections out

	if deps.ShapeType == "county" {
		return deps.syncCounties(shapes)
	} else if deps.ShapeType == "marine" {
		return deps.syncMarineZones(shapes)
	}

	return nil
}

func (deps *SyncShapeData) syncCounties(shapes []domain.DomainShape) error {

	repository := persistence.MongoCountiesRepository{
		Client: deps.MongoClient,
	}

	err := repository.DeleteAll()
	if err != nil {
		return err
	}

	counties := sliceutils.MapFunc(shapes, shapefiles.MapDomainShapeToCounty)

	err = repository.InsertMany(counties)
	if err != nil {
		return err
	}

	return nil
}

func (deps *SyncShapeData) syncMarineZones(shapes []domain.DomainShape) error {

	repostory := persistence.MongoMarineZonesRepository{
		Client: deps.MongoClient,
	}

	err := repostory.DeleteAll()
	if err != nil {
		return err
	}

	marineZones := sliceutils.MapFunc(shapes, shapefiles.MapDomainShapeToMarineZone)

	for _, marineZone := range marineZones {
		err = repostory.InsertMany([]domain.MarineZone{marineZone})
		if err != nil {
			// TODO: need to fix this.  The first MarineZone has a lot of polygons in it
			logging.Error.Println("Could not import: ", err, marineZone.Name, len(marineZone.Geometry.Coordinates[0][0]), len(marineZone.Geometry.Coordinates[0][0][0]), len(marineZone.Geometry.Coordinates[0]))
		} else {
			logging.Debug.Println("Importing: ", marineZone.Name)
		}
	}

	return nil
}

func (deps *SyncShapeData) SyncDomain() error {

	logging.Info.Println("Synchronizing shapefile...")
	logging.Debug.Println("Getting shapes...")
	domainEntities, err := deps.getShapefileEntries()
	if err != nil {
		return err
	}
	logging.Debug.Println("Found ", len(domainEntities), " shapes.")
	logging.Debug.Println("Synchronizing shapes to the database.")

	err = deps.syncDomainTypesInDatabase(domainEntities)
	if err != nil {
		return err
	}

	logging.Debug.Println("Synchronized shapes.")
	logging.Info.Println("Done synchronizing shapes.")

	return nil
}
