package models

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModelProvider struct {
	Model          interface{}
	CollectionName string
	Timestamps     bool
}

var (
	models = map[reflect.Type]ModelProvider{
		reflect.TypeOf(Profile{}): {
			Model:          Profile{},
			CollectionName: "profiles",
			Timestamps:     true,
		},
		reflect.TypeOf(Gender{}): {
			Model:          Gender{},
			CollectionName: "genders",
			Timestamps:     true,
		},
		reflect.TypeOf(Prompt{}): {
			Model:          Prompt{},
			CollectionName: "prompts",
			Timestamps:     true,
		},
	}
)

func GetModels() []ModelProvider {
	var currModels []ModelProvider
	for _, model := range models {
		currModels = append(currModels, model)
	}
	return currModels
}

func beforeCreate(model interface{}) map[string]interface{} {
	timestamps := models[reflect.TypeOf(model)].Timestamps
	if !timestamps {
		return map[string]interface{}{}
	}

	createdAt := primitive.NewDateTimeFromTime(time.Now())
	updatedAt := primitive.NewDateTimeFromTime(time.Now())

	return map[string]interface{}{
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}
}
func Create(ctx *context.Context, db *mongo.Database, model interface{}) (*mongo.InsertOneResult, error) {
	modelProvider := models[reflect.TypeOf(model)]
	collectionName := modelProvider.CollectionName
	additionalFields := beforeCreate(model)

	// Prepare the data to be inserted
	var data map[string]interface{}
	inrec, _ := bson.Marshal(model)
	bson.Unmarshal(inrec, &data)

	// Use reflection to convert the model to a map (optional based on your implementation)
	// For this example, we assume the model is already a map

	for k, v := range additionalFields {
		data[k] = v // Add additional fields
	}
	log.Default().Printf("collection name %s", collectionName)
	collection := db.Collection(collectionName)
	result, err := collection.InsertOne(*ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// beforeUpdate prepares the `updatedAt` field if timestamps are enabled.
func beforeUpdate(model interface{}) map[string]interface{} {
	timestamps := models[reflect.TypeOf(model)].Timestamps
	if !timestamps {
		return map[string]interface{}{}
	}

	updatedAt := primitive.NewDateTimeFromTime(time.Now())
	return map[string]interface{}{
		"updatedAt": updatedAt,
	}
}

// UpdateById updates a document by its ID, setting `updatedAt` if timestamps are enabled.
func UpdateById(ctx *context.Context, db *mongo.Database, model interface{}, id primitive.ObjectID, updateData map[string]interface{}) (*mongo.UpdateResult, error) {
	modelProvider := models[reflect.TypeOf(model)]
	collectionName := modelProvider.CollectionName
	additionalFields := beforeUpdate(model)

	// Merge the `updateData` with `additionalFields` if `timestamps` are enabled
	for k, v := range additionalFields {
		updateData[k] = v
	}

	collection := db.Collection(collectionName)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateOne(*ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateOne updates the first document that matches the filter, setting `updatedAt` if timestamps are enabled.
func UpdateOne(ctx *context.Context, db *mongo.Database, model interface{}, filter bson.M, updateData map[string]interface{}) (*mongo.UpdateResult, error) {
	modelProvider := models[reflect.TypeOf(model)]
	collectionName := modelProvider.CollectionName
	additionalFields := beforeUpdate(model)

	// Merge the `updateData` with `additionalFields`
	for k, v := range additionalFields {
		updateData[k] = v
	}

	collection := db.Collection(collectionName)
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateOne(*ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateMany updates all documents that match the filter, setting `updatedAt` if timestamps are enabled.
func UpdateMany(ctx *context.Context, db *mongo.Database, model interface{}, filter bson.M, updateData map[string]interface{}) (*mongo.UpdateResult, error) {
	modelProvider := models[reflect.TypeOf(model)]
	collectionName := modelProvider.CollectionName
	additionalFields := beforeUpdate(model)

	// Merge the `updateData` with `additionalFields`
	for k, v := range additionalFields {
		updateData[k] = v
	}

	collection := db.Collection(collectionName)
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateMany(*ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOne(ctx *context.Context, db *mongo.Database, model interface{}) *mongo.SingleResult {
	collection := db.Collection(models[reflect.TypeOf(model)].CollectionName)
	var filter map[string]interface{}
	inrec, _ := bson.Marshal(model)
	bson.Unmarshal(inrec, &filter)
	return collection.FindOne(*ctx, filter)
}

func Upsert(
	ctx *context.Context,
	db *mongo.Database,
	where bson.M,
	model interface{},
) (*mongo.UpdateResult, error) {
	collection := db.Collection(models[reflect.TypeOf(model)].CollectionName)

	// Prepare the data to be upserted
	var data map[string]interface{}
	inrec, _ := bson.Marshal(model)
	bson.Unmarshal(inrec, &data)

	// Check if a document matching `where` exists
	var existingData bson.M
	err := collection.FindOne(*ctx, where).Decode(&existingData)
	if err == mongo.ErrNoDocuments {
		additionalFields := beforeCreate(model)
		// Merge `upsertData` and `additionalFields`
		for k, v := range additionalFields {
			data[k] = v
		}

		// Document doesn't exist, insert a new one
		insertResult, err := collection.InsertOne(*ctx, data)
		if err != nil {
			return nil, fmt.Errorf("failed to insert document: %v", err)
		}

		// Return a successful result with UpsertedID populated
		return &mongo.UpdateResult{
			MatchedCount:  0,
			ModifiedCount: 0,
			UpsertedCount: 1,
			UpsertedID:    insertResult.InsertedID,
		}, nil
	} else if err != nil {
		// Other errors during FindOne
		return nil, fmt.Errorf("failed to find document: %v", err)
	}

	additionalFields := beforeUpdate(model)
	// Merge `upsertData` and `additionalFields`
	for k, v := range additionalFields {
		data[k] = v
	}

	// Document exists, update it
	filter := bson.M{"_id": existingData["_id"]} // Use the existing document's ID
	update := bson.M{"$set": data}               // Use `$set` for partial updates

	updateResult, err := collection.UpdateOne(*ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}

	return updateResult, nil
}

func Find(ctx *context.Context, db *mongo.Database, model interface{}, opts *options.FindOptions) (*mongo.Cursor, error) {
	collection := db.Collection(models[reflect.TypeOf(model)].CollectionName)
	var filter map[string]interface{}
	inrec, _ := bson.Marshal(model)
	bson.Unmarshal(inrec, &filter)
	return collection.Find(*ctx, filter, opts)
}

func CountAllAndFind(ctx *context.Context, db *mongo.Database, model interface{}, opts *options.FindOptions) (*int64, *mongo.Cursor, error) {
	// Get collection based on the model type
	collection := db.Collection(models[reflect.TypeOf(model)].CollectionName)

	// Serialize the model to BSON
	var filter map[string]interface{}
	inrec, err := bson.Marshal(model)
	if err != nil {
		return nil, nil, err
	}
	err = bson.Unmarshal(inrec, &filter)
	if err != nil {
		return nil, nil, err
	}

	// Count documents
	count, err := collection.CountDocuments(*ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	// Find matching documents
	cursor, err := collection.Find(*ctx, filter, opts)
	if err != nil {
		return nil, nil, err
	}

	return &count, cursor, nil
}
