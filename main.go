package main

import(
	"fmt"
	"encoding/json"
	"time"
	"net/http"
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
	


)

type Users struct{
	Id primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	DOB string `json:"dob" bson:"dob"`
	Phone string `json:"ph,omitempty" bson:"ph,omitempty"`
	email string `json:"email,omitempty" bson:"email,omitempty"`
	Creation time.Time`json:"created,omitempty" bson:"created,omitempty"`
}

type Contact struct{
	userId_one  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	userId_two  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	contact_time time.Time `json:"cont_time,omitempty" bson:"cont_time,omitempty"`
}

var client *mongo.Client

func CreateUser(response http.ResponseWriter, request *http.Request)  {
	response.Header().Add("content-type", "application/json")
	var users Users
	json.NewDecoder(request.Body).Decode(&users)
	ctx, _ :=context.WithTimeout(context.Background(), 10*time.Second)
	collection:= client.Database("mydatabase").Collection("Users")
	result, _ :=collection.InsertOne(ctx, users)
	json.NewEncoder(response).Encode(result)
}
func GetUser(response http.ResponseWriter, request *http.Request)  {
	response.Header().Add("content-type", "application/json")
	params:=mux.Vars(request)
	id,_ :=primitive.ObjectIDFromHex(params["Id"])
	var users Users
	collection:= client.Database("mydatabase").Collection("Users")
	ctx, _ :=context.WithTimeout(context.Background(), 10*time.Second)
	err:=collection.FindOne(ctx, Users{Id: id}).Decode(&users)
	if err!=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([] byte(`{"message": "` +err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(users)
}

func AddContact(response http.ResponseWriter, request *http.Request)  {
	response.Header().Add("content-type", "application/json")
	var contacts Contact
	json.NewDecoder(request.Body).Decode(&contacts)
	ctx, _ :=context.WithTimeout(context.Background(), 10*time.Second)
	collection:= client.Database("mydatabase").Collection("Contact")
	result, _:=collection.InsertOne(ctx,contacts)
	json.NewEncoder(response).Encode(result)
}

func main()  {
	fmt.Println("Application is starting...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://datastore:27017"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
	router :=mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/id", GetUser).Methods("GET")
	router.HandleFunc("/contacts/{id}", AddContact).Methods("POST")
	http.ListenAndServe(":12345", router)


}
