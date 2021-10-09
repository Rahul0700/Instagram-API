package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	/*
		:param client: <mongo.Client> Identify the associated resource
		:param ctx: <context.Context> Allows to set deadlines for the Disconnect process
		:param cancel: <context.CancelFunc> Cancels the context
	*/
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	/*
		:param uri: <string> Resource Identifier of the mongodb instance hosted in the cloud
		:return client: <mongo.Client> Will be used to perform further database operations
		:return ctx: <context.Context> Allows to set deadlines for the process
		:return cancel: <context.CancelFunc> Allows to cancel the associated context
		:return err: <error> Returns nil if connection successful
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func insertDocument(dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	/*
		:param dataBase: <string> Name of the database where the document should be inserted
		:param col: <string> Collection name where the document is supposed to be installed
		:param doc: <interface{}> Holds the document to be inserted
		:return result: <*mongo.InsertOneResult> Contains the particular instance of the created document
		:return err: <error> Returns nil if connection successful
	*/
	client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	if err != nil {
		panic(err)
		fmt.Print(cancel)
	}
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func getDocument(dataBase, col string, id string) (primitive.M, error) {
	/*
		:param dataBase: <string> Name of the database where the document should be searched at
		:param col: <string> Collection name where the document is supposed to be searched at
		:param id: <string> The id of the document the user is looking for
		:return result: <primitive.M> The retrieved document instance
		:return err: <error> Returns nil if connection successful
	*/
	client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	if err != nil {
		panic(err)
		fmt.Print(cancel)
	}
	var result bson.M
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	collection := client.Database(dataBase).Collection(col)
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
	return result, err
}
