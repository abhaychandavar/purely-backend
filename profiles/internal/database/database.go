package database

import (
	"context"
	"fmt"
	"log"
	"profiles/internal/config"
	"profiles/internal/database/models"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	Db() *mongo.Database
	Client() *mongo.Client
	init()
}

type service struct {
	db     *mongo.Database
	client *mongo.Client
}

var (
	instance Service
	once     sync.Once
)

func createIndexes(client *mongo.Client, model interface{}, collectionName string) error {
	collection := client.Database(config.GetConfig().Db).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", t.Kind())
	}

	var indexModels []mongo.IndexModel

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		bsonTag := field.Tag.Get("bson")
		uniqueTag := field.Tag.Get("unique")

		// If the field has a BSON tag and is unique
		if bsonTag != "" && uniqueTag == "true" {
			parts := strings.Split(bsonTag, ",")
			tagName := parts[0]
			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: tagName, Value: 1}}, // Create an ascending index
				Options: options.Index().SetUnique(true),
			}
			indexModels = append(indexModels, indexModel)
		}
	}

	// Create the indexes in MongoDB
	if len(indexModels) > 0 {
		if _, err := collection.Indexes().CreateMany(ctx, indexModels); err != nil {
			return fmt.Errorf("failed to create indexes: %v", err)
		}
	}

	return nil
}

func (s *service) init() {
	models := models.GetModels()
	for _, modelProvider := range models {
		if err := createIndexes(s.client, modelProvider.Model, modelProvider.CollectionName); err != nil {
			log.Printf("error creating indexes for model %T: %v", modelProvider.Model, err)
		}
	}
}

func Mongo() Service {
	once.Do(func() {
		log.Default().Printf("Connecting to mongo")
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.GetConfig().MongoConnUrl))
		if err != nil {
			log.Fatal(err)
		}
		instance = &service{
			client: client,
			db:     client.Database(config.GetConfig().Db),
		}
		instance.init()
	})
	return instance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.client.Ping(ctx, nil); err != nil {
		log.Printf("db down: %v", err)
		return map[string]string{"message": "db down"}
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) Client() *mongo.Client {
	return s.client
}

func (s *service) Db() *mongo.Database {
	return s.db
}
