package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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

func getMD5Hash(text string) string {
	/*
		:param text: <string> Pasasword to be encrypted
		:return hex: <string> Encrypted hash of the password
	*/
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func userRequestHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

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
		createuser, err := insertDocument("Instagram-API", "Users", userBson)
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
		result, err := getDocument("Instagram-API", "Users", id)

		// Response
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
