package sql

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type lend struct {
	BookName   string
	Index      string
	Name       string
	StudentId  string
	Academy    string
	StartDate  string
	ReturnDate string
	Status     string
	Seashell   string
}

func Addlend(lend map[string]string) {

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

	collection := client.Database(database).Collection("lends")

	insertOne, err := collection.InsertOne(ctx, lend)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Document: ", insertOne.InsertedID)
}

func Findlend() []*lend {

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
	collection := client.Database(database).Collection("lends")

	findOptions := options.Find()

	var results []*lend

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var result lend
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

func ReturnId(studentid string, index string) {

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
	collection := client.Database(database).Collection("lends")

	filter := bson.D{{"studentid", studentid}, {"index", index}}
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
