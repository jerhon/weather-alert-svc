package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/sliceutils"
)

type MarineZonesRepository interface {
	DeleteAll() error
	InsertMany(counties []domain.MarineZone) error
}

type MongoMarineZonesRepository struct {
	Client *mongo.Client
}

func (dep *MongoMarineZonesRepository) DeleteAll() error {
	collection := dep.getCollection()
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	return nil
}

func (dep *MongoMarineZonesRepository) InsertMany(marineZones []domain.MarineZone) error {
	collection := dep.getCollection()
	mongoCounties := sliceutils.MapFunc(marineZones, func(county domain.MarineZone) interface{} {
		return county
	})
	_, err := collection.InsertMany(context.TODO(), mongoCounties)
	if err != nil {
		return err
	}
	return nil
}

func (dep *MongoMarineZonesRepository) getCollection() *mongo.Collection {
	return dep.Client.Database("weather-alerts").Collection("marinezones")
}
