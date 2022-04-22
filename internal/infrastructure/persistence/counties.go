package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/sliceutils"
)

type CountiesRepository interface {
	DeleteAll() error
	InsertMany(counties []domain.County) error
	FindCountyBySame(same string) error
}

type MongoCountiesRepository struct {
	Client *mongo.Client
}

func (repo *MongoCountiesRepository) DeleteAll() error {
	collection := repo.getCountiesCollection()
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoCountiesRepository) InsertMany(counties []domain.County) error {
	collection := repo.getCountiesCollection()
	mongoCounties := sliceutils.MapFunc(counties, func(alert domain.County) interface{} {
		return alert
	})
	_, err := collection.InsertMany(context.TODO(), mongoCounties)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoCountiesRepository) FindCountyBySame(same string) (*domain.County, error) {
	collection := repo.getCountiesCollection()
	result := collection.FindOne(context.TODO(), bson.D{{"same", same}}, options.FindOne())
	if result.Err() != nil {
		return nil, result.Err()
	}

	county := &domain.County{}
	err := result.Decode(county)
	if err != nil {
		return nil, err
	}
	return county, nil
}

func (repo *MongoCountiesRepository) getCountiesCollection() *mongo.Collection {
	return repo.Client.Database("weather-alerts").Collection("counties")
}
