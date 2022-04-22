package usecases

import (
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/internal/infrastructure/nwsshapefiles"
	"weather-alerts-service/internal/infrastructure/persistence"
)

type ImportCountyDeps struct {
	ShapefileReader *nwsshapefiles.CountyReader
	Repository      *persistence.MongoCountiesRepository
}

func (deps *ImportCountyDeps) getShapefileCounties() ([]domain.County, error) {

	counties, err := deps.ShapefileReader.GetAllCounties()
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

	counties, err := deps.getShapefileCounties()
	if err != nil {
		return err
	}

	err = deps.syncCounties(counties)
	if err != nil {
		return err
	}

	return nil
}
