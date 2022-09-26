package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB ...

func ConnectDB () *mongo.Client {
	godotenv.Load()
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}	

	fmt.Println("Connected to MongoDB")

	return client
}

func GetCollection (client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(os.Getenv("DB_NAME")).Collection(collectionName);
	return collection
}

func GetDB() *mongo.Client {
	return DB
}

var DB *mongo.Client = ConnectDB()