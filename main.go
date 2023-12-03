package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Person struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user:pwd@localhost:27017"))

	if err != nil {
		panic(err.Error())
	}

	//链接database
	db := client.Database("mydb")
	//读取集合，没有则创建
	collection := db.Collection("persons")

	// 创建文档
	person := Person{Name: "John Doe", Age: 30}
	_, err = collection.InsertOne(context.Background(), person)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created document with ID:", person.ID)

	// 读取文档
	filter := bson.M{"name": "John Doe"}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	var result Person
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
		fmt.Println("Read document:", result)
	}

	if err = cursor.Err(); err != nil {
		panic(err)
	}

	// 更新文档
	filter = bson.M{"name": "John Doe"}
	update := bson.M{"$set": bson.M{"age": 35}}
	updateResult, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated", updateResult.ModifiedCount, "document")

	//获取集合名
	collectionNames, err := db.ListCollections(ctx, bson.M{})
	if err != nil {
		panic(err.Error())
		return
	}

	fmt.Println("集合名：", collectionNames)

	//断开链接
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err.Error())
		}
	}()

}
