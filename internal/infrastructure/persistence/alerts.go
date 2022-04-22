package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/sliceutils"
)

// TODO does returning the Client take a connection from the pool, or is that only allocated when something is called?

type AlertRepository interface {
	FindExistingAlerts(alerts []domain.Alert) ([]domain.Alert, error)
	InsertAlerts(alert []domain.Alert) error
}

type MongoAlertRepository struct {
	Client *mongo.Client
}

func NewAlertMongoRepository(mongoClient *mongo.Client) AlertRepository {
	repository := new(MongoAlertRepository)
	repository.Client = mongoClient
	return repository
}

// FindExistingAlerts finds any alerts that already exist in the alerts collection.
func (repo *MongoAlertRepository) FindExistingAlerts(alerts []domain.Alert) ([]domain.Alert, error) {
	var ret []domain.Alert
	opts := options.Find().SetProjection(bson.D{{"originid", 1}})
	alertIds := sliceutils.MapFunc(alerts, GetAlertOriginId)
	alertsCollection := repo.getAlertsCollection()
	cursor, err := alertsCollection.Find(context.TODO(), bson.D{{"originid", bson.D{{"$in", alertIds}}}}, opts)
	if err != nil {
		return nil, err
	}
	var allAlerts []domain.Alert
	err = cursor.All(context.TODO(), &allAlerts)
	if err != nil {
		return nil, err
	}

	for _, alert := range allAlerts {
		idx := slices.IndexFunc(alerts, func(e domain.Alert) bool {
			return e.OriginId == alert.OriginId
		})
		if idx > 0 {
			ret = append(ret, alert)
		}
	}
	return ret, nil
}

func GetAlertOriginId(alert domain.Alert) string {
	return alert.OriginId
}

func (repo *MongoAlertRepository) InsertAlerts(alerts []domain.Alert) error {
	alertsCollection := repo.getAlertsCollection()

	// This is a really unfortunate consequence of the go language, hoping it will get better once the go driver implements generics
	mongoAlerts := sliceutils.MapFunc(alerts, func(alert domain.Alert) interface{} {
		return alert
	})

	_, err := alertsCollection.InsertMany(context.TODO(), mongoAlerts)
	return err
}

func (repo *MongoAlertRepository) getAlertsCollection() *mongo.Collection {
	return repo.Client.Database("weather-alerts").Collection("alerts")
}
