package sql

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type time struct {
	Username string
	Time     string
}

func Findtime(username string) []*time {

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

	collection := client.Database(database).Collection("times")

	var results []*time

	sort := bson.D{
		bson.E{"time", -1},
	}

	opts := &options.FindOptions{
		Sort: sort,
	}

	cur, err := collection.Find(context.TODO(), bson.D{{"username", username}}, opts)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var result time
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

func Addtime(username string, mytime string) {

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

	collection := client.Database(database).Collection("times")

	time := time{Username: username, Time: mytime}

	insertOne, err := collection.InsertOne(ctx, time)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Document: ", insertOne.InsertedID)
}
