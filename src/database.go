package main

import (
	"context"
	"time"

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

func main() {

	client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	if err != nil {
		panic(err)
	}

	// Called when the main is returned
	defer close(client, ctx, cancel)
}
