package usecases

import (
	"time"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/internal/infrastructure/persistence"
	"weather-alerts-service/internal/infrastructure/weatherservice"
	"weather-alerts-service/internal/logging"
	"weather-alerts-service/pkg/sliceutils"
)

type SyncAlertDependencies struct {
	AlertSource         weatherservice.AlertSource
	AlertRepository     persistence.AlertRepository
	ImportLogRepository persistence.ImportLogRepository
}

func (deps SyncAlertDependencies) SyncAlerts() (int, error) {

	lastImport, err := deps.ImportLogRepository.GetLastImport("active-alerts")
	if err != nil {
		logging.Error.Println("Cannot retrieve last import record", err)
	}

	ifModifiedSince := ""
	if lastImport != nil {
		ifModifiedSince = lastImport.LastModified
	}

	alerts, lastModified, err := deps.AlertSource.GetActiveAlerts(ifModifiedSince)
	if err != nil {
		return 0, err
	}

	if len(alerts) > 0 {

		newAlerts, err := deps.getNewAlerts(alerts)
		if err != nil {
			return 0, err
		}

		err = deps.AlertRepository.InsertAlerts(newAlerts)
		if err != nil {
			return 0, err
		}

		err = deps.ImportLogRepository.Insert(domain.ImportLog{
			Type:           "active-alerts",
			LastModified:   lastModified,
			ImportedTime:   time.Now().UTC(),
			ImportedAlerts: len(newAlerts),
		})
		if err != nil {
			return 0, err
		}

		return len(newAlerts), err
	}
	return 0, nil
}

func (deps SyncAlertDependencies) getNewAlerts(alerts []domain.Alert) ([]domain.Alert, error) {

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
