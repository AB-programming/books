package sql

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type book struct {
	Name     string
	Index    string
	Classify string
	Type     string
	Num      int32
}

type conf struct {
	Url      string `yaml:"url"`
	Database string `yaml:"database"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("./yaml/conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

func Find() []*book {

	var c conf
	conf := c.getConf()
	url := conf.Url
	database := conf.Database

	clientOptions := options.Client().ApplyURI(url)
	var ctx = context.TODO()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接MongoDB成功!")
	defer client.Disconnect(ctx)

	collection := client.Database(database).Collection("books")
	findOptions := options.Find()

	var results []*book

	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var result book
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

func Updatebook(index string) {
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

	collection := client.Database(database).Collection("books")

	var result book
	filter := bson.D{{"index", index}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	num := result.Num - 1

	filter = bson.D{{"index", index}}

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", bson.D{
			{"num", num}},
		}}
	result2, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	if result2.MatchedCount != 0 {
		fmt.Printf("Matched %v documents and updated %v documents.\n", result2.MatchedCount, result2.ModifiedCount)
	}
	if result2.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result2.UpsertedID)
	}
}

func Findbook(classify string, mytype string, name string) []*book {

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

	collection := client.Database(database).Collection("books")

	findOptions := options.Find()

	var results []*book

	var bs bson.D
	if classify == "" && mytype == "" && name == "" {
		bs = bson.D{{}}
	} else if classify == "" && mytype == "" {
		bs = bson.D{{"name", name}}
	} else if mytype == "" && name == "" {
		bs = bson.D{{"classify", classify}}
	} else if classify == "" && name == "" {
		bs = bson.D{{"type", mytype}}
	} else if classify == "" {
		bs = bson.D{{"type", mytype}, {"name", name}}
	} else if mytype == "" {
		bs = bson.D{{"classify", classify}, {"name", name}}
	} else if name == "" {
		bs = bson.D{{"classify", classify}, {"type", mytype}}
	} else {
		bs = bson.D{{"classify", classify}, {"type", mytype}, {"name", name}}
	}

	cur, err := collection.Find(context.TODO(), bs, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var result book
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

func Addbook(bookname string, index string, classify string, mytype string, num int32) {

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

	collection := client.Database(database).Collection("books")

	book := book{Name: bookname, Index: index, Classify: classify, Type: mytype, Num: num}
	insertOne, err := collection.InsertOne(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Document: ", insertOne.InsertedID)
}

func Deletebook(index string) {

	var c conf
	conf := c.getConf()
	url := conf.Url
	database := conf.Database

	clientOptions := options.Client().ApplyURI(url)
	var ctx = context.TODO()
	// Connect to MongoDB
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

	collection := client.Database(database).Collection("books")

	filter := bson.D{{"index", index}}
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
