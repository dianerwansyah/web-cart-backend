package helper

import (
	"context"
	"log"
	"reflect"

	"github.com/dianerwansyah/web-cart-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitDB(uri string, collections []string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	cfg := GetConfig()
	db := client.Database(cfg.Server.MongoDB)
	for _, collection := range collections {
		CreateCollectionIfNotExists(db, collection)
	}

	return client
}

func CreateCollectionIfNotExists(db *mongo.Database, collectionName string) {
	collections, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatalf("Error listing collections: %v", err)
	}

	collectionExists := false
	for _, name := range collections {
		if name == collectionName {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		err := db.CreateCollection(context.Background(), collectionName)
		if err != nil {
			log.Fatalf("Error creating collection: %v", err)
		}
		log.Printf("Collection %s created.", collectionName)
	}
}

func GetDBClient() *mongo.Client {
	return client
}

func SetDBClient(mongoClient *mongo.Client) {
	client = mongoClient
}

func GetAllModels() []interface{} {
	return []interface{}{
		model.Category{},
		model.Product{},
		model.Cart{},
		model.Coupon{},
		model.History{},
	}
}

func GetTableNames(models []interface{}) []string {
	var tableNames []string
	for _, model := range models {
		tableNames = append(tableNames, getTableName(model))
	}
	return tableNames
}

func getTableName(model interface{}) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	method, found := modelType.MethodByName("TableName")
	if found {
		results := method.Func.Call([]reflect.Value{reflect.ValueOf(model)})
		return results[0].String()
	}
	log.Fatalf("getTableName method not found for model %s", modelType.Name())
	return ""
}

func GetCollection(TableName string) *mongo.Collection {
	collectionName := TableName

	cfg := GetConfig()
	db := client.Database(cfg.Server.MongoDB)
	return db.Collection(collectionName)
}

func GetTableName(models interface{}) string {
	return getTableName(models)
}
