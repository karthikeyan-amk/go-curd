package initializer

import (
	"context"
	"log"

	"github.com/karthikeyan-amk/go-curd/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TestCollection *mongo.Collection

func Database() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://karthikeyan:Ik1JW84ClCJkljzO@cluster0.mewwxmt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	TestCollection = client.Database("test").Collection("test")
	log.Println("Connected to MongoDB!")
}

func CreateMovie(user *models.Movie) error {
	// Insert user into the database
	_, err := TestCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	log.Println("User created successfully!")
	return nil
}
