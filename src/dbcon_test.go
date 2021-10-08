package main

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ping(client *mongo.Client, ctx context.Context) error {
	/*
		:param client: <mongo.Client> Ping the mongodb client
		:param ctx: <context.Context> Allows to set deadlines for Ping
		:return err: <error> Returns nil if connection successful
	*/

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func dbConnTest(t *testing.T) {
	client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	got := ping(client, ctx)
	if got != nil {
		t.Errorf("Database Connection test: FAIL")
	}
}
