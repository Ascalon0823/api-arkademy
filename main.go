package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Sale struct {
	Item     string  `bson:"item"`
	Price    float64 `bson:"price"`
	Quantity int     `bson:"quantity"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Environment")
	}
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI_LOCAL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	var results []Sale
	coll := client.Database("mongodbVSCodePlaygroundDB").Collection("sales")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
