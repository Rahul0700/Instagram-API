package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username string
	Password string
	Email    string
}

type UserPostResponse struct {
	Status string
	Id     interface{}
}

type UserGetResponse struct {
	Status   string
	Id       string
	Username string
	Password string
	Email    string
}

func getMD5Hash(text string) string {
	/*
		:param text: <string> Pasasword to be encrypted
		:return hex: <string> Encrypted hash of the password
	*/
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func userRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// POST req to handle user registration
	case "POST":
		// Decode Json request body to struct layout
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			panic(err)
		}

		// Password Encryption using MD5 Algo
		rawpassword := user.Password
		user.Password = getMD5Hash(rawpassword)
		userBson, err := bson.Marshal(user)
		if err != nil {
			panic(err)
		}

		//Insert Document
		client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
		if err != nil {
			panic(err)
			fmt.Print(cancel)
		}
		createuser, err := insertDocument(client, ctx, "Instagram-API", "Users", userBson)
		if err != nil {
			panic(err)
		}

		//Response
		userresponse := UserPostResponse{}
		userresponse.Id = createuser.InsertedID
		userresponse.Status = "success"
		userJson, err := json.Marshal(userresponse)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(userJson)

	// GET req to get user details using id
	case "GET":
		// Parse userid from url
		id := r.URL.Path[7:]

		//Find doocument using id filter
		client, ctx, cancel, err := connect("mongodb+srv://rahul:QrpiHbW1srNcm9I5@cluster0.aumtt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
		if err != nil {
			panic(err)
			fmt.Print(cancel)
		}
		var result UserGetResponse
		collection := client.Database("Instagram-API").Collection("Users")
		docID, err := primitive.ObjectIDFromHex(id)
		err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
		if err != nil {
			panic(err)
		}

		// Response
		result.Id = id
		result.Status = "success"
		userJson, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"failure", "Message": "400 Bad Request"}`))
	}
}
