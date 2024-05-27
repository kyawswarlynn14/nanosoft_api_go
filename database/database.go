package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}
	fmt.Println("Successfully Connected to the mongodb")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Nanosoft").Collection(CollectionName)
	return collection

}

func ServiceData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Nanosoft").Collection(CollectionName)
	return productcollection
}

func ProjectData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Nanosoft").Collection(CollectionName)
	return productcollection
}

func RemarkData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Nanosoft").Collection(CollectionName)
	return productcollection
}

func EmailData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Nanosoft").Collection(CollectionName)
	return productcollection
}

func BlogData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Nanosoft").Collection(CollectionName)
	return productcollection
}
