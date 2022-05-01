package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(host string, user string, password string) (*mongo.Client, error) {
	var clientOptions = options.Client().SetMaxPoolSize(20)
	if len(host) > 0 {
		clientOptions = clientOptions.SetHosts([]string{host})
	}
	if len(user) > 0 && len(password) > 0 {
		clientOptions = clientOptions.SetAuth(options.Credential{Username: user, Password: password})
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)
	return client, err
}

type SyncDomainTypeRepository[T any] interface {
	DeleteAll() error
	InsertMany(domainEntries []T) error
}
