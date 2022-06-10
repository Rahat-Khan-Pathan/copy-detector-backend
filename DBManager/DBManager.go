package DBManager

import (
	"context"
	"log"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var configErr = godotenv.Load()
var dbURL string = os.Getenv("DB_SOURCE_URL")

var SystemCollections VAICollections

type VAICollections struct {
	Users   *mongo.Collection
	Courses *mongo.Collection
	Exams   *mongo.Collection
	Codes   *mongo.Collection
}

func getMongoDbConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := getMongoDbConnection()
	if err != nil {
		return nil, err
	}
	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

func InitCollections() bool {
	if configErr != nil {
		return false
	}
	var err error

	SystemCollections.Users, err = GetMongoDbCollection("copy-detector", "users")
	if err != nil {
		return false
	}
	SystemCollections.Courses, err = GetMongoDbCollection("copy-detector", "courses")
	if err != nil {
		return false
	}
	SystemCollections.Exams, err = GetMongoDbCollection("copy-detector", "exams")
	if err != nil {
		return false
	}
	SystemCollections.Codes, err = GetMongoDbCollection("copy-detector", "codes")
	if err != nil {
		return false
	}
	return true
}
