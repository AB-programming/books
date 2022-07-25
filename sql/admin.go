package sql

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type admin struct {
	Username string
	Password string
	Nickname string
	Sex      string
	Age      string
}

func FindAdmin() []*admin {

	var c conf
	conf := c.getConf()
	url := conf.Url
	database := conf.Database

	clientOptions := options.Client().ApplyURI(url)
	var ctx = context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接MongoDB成功!")
	defer client.Disconnect(ctx)

	collection := client.Database(database).Collection("admins")

	findOptions := options.Find()

	var results []*admin

	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var result admin
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return results
}

func FindOneAdmin(username string) *admin {

	var c conf
	conf := c.getConf()
	url := conf.Url
	database := conf.Database

	clientOptions := options.Client().ApplyURI(url)
	var ctx = context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接MongoDB成功!")
	defer client.Disconnect(ctx)

	collection := client.Database(database).Collection("admins")

	var result *admin
	filter := bson.D{{"username", username}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	return result
}
