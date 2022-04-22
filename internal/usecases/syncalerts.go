package usecases

import (
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/infrastructure/weatherservice"
	"weather-alerts-service/pkg/sliceutils"
)

type SyncAlertDependencies struct {
	AlertSource      weatherservice.AlertSource
	AlertRepository  persistence.AlertRepository
	CountyRepository persistence.CountiesRepository
}

func (deps SyncAlertDependencies) SyncAlerts() (int, error) {
	alerts, err := deps.AlertSource.GetActiveAlerts()
	if err != nil {
		return 0, err
	}

	newAlerts, err := deps.getNewAlerts(alerts)
	if err != nil {
		return 0, err
	}

	err = deps.AlertRepository.InsertAlerts(newAlerts)
	if err != nil {
		return 0, err
	}

	return len(newAlerts), err
}

func (deps SyncAlertDependencies) populateAlertGeometry(alerts []domain.Alert) {
	for _, alert := range alerts {
		if alert.Geometry != nil {

		}
	}
}

func (deps SyncAlertDependencies) getNewAlerts(alerts []domain.Alert) ([]domain.Alert, error) {

	// Find any existing alerts
	ret := make([]domain.Alert, 0)
	oldAlerts, err := deps.AlertRepository.FindExistingAlerts(alerts)
	if err != nil {
		return nil, err
	}
	oldAlertMap := sliceutils.ToMapFunc(oldAlerts, func(a domain.Alert) string {
		return a.OriginId
	})

	for _, alert := range alerts {
		_, isPresent := oldAlertMap[alert.OriginId]
		if !isPresent {
			ret = append(ret, alert)
		}
	}

	return ret, nil
}
