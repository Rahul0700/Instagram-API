package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Posts struct {
	Userid    string
	Caption   string
	Imageurl  string
	Timestamp time.Time
}

type PostsPostResponse struct {
	Status string
	Id     interface{}
}

func postsRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// POST req to handle Post creation
	case "POST":
		// Decode Json request body to struct layout
		posts := Posts{}
		err := json.NewDecoder(r.Body).Decode(&posts)
		if err != nil {
			panic(err)
		}

		posts.Timestamp = time.Now()
		postsBson, err := bson.Marshal(posts)
		if err != nil {
			panic(err)
		}

		//Insert Document
		client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
		if err != nil {
			panic(err)
			fmt.Print(cancel)
		}
		createposts, err := insertDocument(client, ctx, "Instagram-API", "Posts", postsBson)
		if err != nil {
			panic(err)
		}

		//Response
		postsresponse := PostsPostResponse{}
		postsresponse.Id = createposts.InsertedID
		postsresponse.Status = "success"
		postsJson, err := json.Marshal(postsresponse)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(postsJson)

	// GET req to get post details using id
	case "GET":
		// Parse userid from url
		id := r.URL.Path[7:]

		//Find doocument using id filter
		client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
		if err != nil {
			panic(err)
			fmt.Print(cancel)
		}
		result, err := getDocument(client, ctx, "Instagram-API", "Posts", id)

		// Response
		postsJson, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(postsJson)

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"failure", "Message": "400 Bad Request"}`))
	}
}
