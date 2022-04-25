package usecases

import (
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/internal/infrastructure/nwsshapefiles"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/logging"
)

type ImportCountyDeps struct {
	ShapefileReader *nwsshapefiles.DomainShapeReader[domain.County]
	Repository      *persistence.MongoCountiesRepository
}

func (deps *ImportCountyDeps) getShapefileCounties() ([]domain.County, error) {

	counties, err := deps.ShapefileReader.GetAll()
	if err != nil {
		return nil, err
	}

	return counties, nil
}

func (deps *ImportCountyDeps) syncCounties(counties []domain.County) error {

	// TODO: make this non-destructuve with an upsert so counties can be updated in real time
	err := deps.Repository.DeleteAll()
	if err != nil {
		return err
	}

	err = deps.Repository.InsertMany(counties)
	if err != nil {
		return err
	}

	return nil
}

func (deps *ImportCountyDeps) SyncCounties() error {

	logging.Info.Println("Synchronizing counties shapefile...")
	logging.Debug.Println("Getting shapes...")
	counties, err := deps.getShapefileCounties()
	if err != nil {
		return err
	}
	logging.Debug.Println("Found ", len(counties), " shapes.")
	logging.Debug.Println("Synchronizing counties to the database.")

	err = deps.syncCounties(counties)
	if err != nil {
		return err
	}

	logging.Debug.Println("Synchronized counties.")
	logging.Info.Println("Done synchronizing counties.")

	return nil
}
