package main

import ("context"
"encoding/json"
"fmt"
"go.mongodb.org/mongo-driver/mongo/options"
"net/http"
"time"
"github.com/gorilla/handlers"
"github.com/gorilla/mux"
"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/bson/primitive"
"go.mongodb.org/mongo-driver/mongo")
type Person struct{
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Color     string             `json:"Color,omitempty" bson:"Color,omitempty"`
	Participation int64          `json:"Participation,omitempty" bson:"Participation,omitempty"`
}
var  client *mongo.Client
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("lega-green").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("lega-green").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}
func main (){
	fmt.Println("Starting the application...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://ramiroJCB:7pMIHlYs6OMzQylH@cluster0.uuj77.mongodb.net/lega-people?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	headers :=handlers.AllowedHeaders([]string{"X-Requested-With","Content-type","Authotization"})
	methods := handlers.AllowedMethods([]string{"GET","POST","PUT","DELETE"})
	origins:=handlers.AllowedOrigins([]string{"*"})
	
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")

	http.ListenAndServe(":8080",handlers.CORS(headers,methods,origins)(router))
}