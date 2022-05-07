package weatherservice

import (
	"time"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/geojson"
	"weather-alerts-service/pkg/nwsclient"
	"weather-alerts-service/pkg/sliceutils"
)

type AlertSource interface {
	GetActiveAlerts(lastModified string) ([]domain.Alert, string, error)
}

type AlertSourceAdapter struct {
	Client *nwsclient.NwsClient
}

func NewAlertsAdapter(applicationString string) *AlertSourceAdapter {
	return &AlertSourceAdapter{
		Client: nwsclient.NewNwsClient(applicationString),
	}
}

func (adapter *AlertSourceAdapter) GetActiveAlerts(lastModified string) ([]domain.Alert, string, error) {
	nwsAlerts, err := adapter.Client.GetActiveAlerts(lastModified)
	if err != nil {
		return nil, "", err
	}
	domainAlerts := mapFeatureCollectionToAlerts(nwsAlerts.Alerts)
	return domainAlerts, nwsAlerts.LastModified, nil
}

// MapToDomain takes an alert from the geojson and maps into it's form to be stored as a domain type
func mapFeatureCollectionToAlerts(featureCollection *geojson.FeatureCollection[nwsclient.AlertProperties]) []domain.Alert {
	return sliceutils.MapFunc(featureCollection.Features, mapFeatureToAlert)
}

func mapFeatureToAlert(item geojson.Feature[nwsclient.AlertProperties]) domain.Alert {

	alert := domain.Alert{
		OriginId:    item.Properties.Id,
		Type:        item.Properties.Event,
		UGC:         item.Properties.Geocode.UGC,
		SAME:        item.Properties.Geocode.SAME,
		AreaDesc:    item.Properties.AreaDesc,
		Event:       item.Properties.Event,
		Category:    item.Properties.Category,
		Certainty:   item.Properties.Certainty,
		Description: item.Properties.Description,
		Ends:        parseTime(item.Properties.Ends),
		Sent:        parseTime(item.Properties.Sent),
		Effective:   parseTime(item.Properties.Effective),
		Onset:       parseTime(item.Properties.Onset),
		Expires:     parseTime(item.Properties.Expires),
		MessageType: item.Properties.MessageType,
		Sender:      item.Properties.Sender,
		Headline:    item.Properties.Headline,
		Instruction: item.Properties.Instruction,
		References:  sliceutils.MapFunc(item.Properties.References, mapAlertReferenceIds),
		Severity:    item.Properties.Severity,
		Status:      item.Properties.Status,
		Urgency:     item.Properties.Urgency,
		Response:    item.Properties.Response,
		SenderName:  item.Properties.SenderName,
	}

	if item.Geometry != nil {
		alert.Geometry = item.Geometry.GetMultiPolygon()
	}

	return alert
}

func mapAlertReferenceIds(references nwsclient.AlertReference) string {
	return references.Identifier
}

func parseTime(timeString string) time.Time {
	if timeString == "" {
		return time.Time{}
	}

	ret, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time.Time{}
	}

	return ret
}
