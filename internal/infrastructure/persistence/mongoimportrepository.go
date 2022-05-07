package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"weather-alerts-service/internal/domain"
)

type ImportLogRepository interface {
	Insert(importProcess domain.ImportLog) error
	GetLastImport(importType string) (*domain.ImportLog, error)
}

type MongoImportLogRepository struct {
	Client *mongo.Client
}

func NewMongoImportLogRepository(client *mongo.Client) MongoImportLogRepository {
	return MongoImportLogRepository{Client: client}
}

func (repo MongoImportLogRepository) Insert(importProcess domain.ImportLog) error {

	collection := repo.getCollection()

	_, err := collection.InsertOne(context.TODO(), importProcess)
	if err != nil {
		return err
	}

	return nil
}

func (repo MongoImportLogRepository) GetLastImport(importType string) (*domain.ImportLog, error) {
	collection := repo.getCollection()

	result := collection.FindOne(context.TODO(), bson.D{{}}, &options.FindOneOptions{
		Sort: bson.D{{"importedtime", -1}},
	})

	if result.Err() != nil {
		return nil, result.Err()
	}

	log := &domain.ImportLog{}
	err := result.Decode(log)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (repo *MongoImportLogRepository) getCollection() *mongo.Collection {
	return repo.Client.Database("weather-alerts").Collection("importlog")
}
