package controllers

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/karthikeyan-amk/go-curd/initializer"
	"github.com/karthikeyan-amk/go-curd/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateMovie(c *gin.Context) {

	var Movie models.Movie

	if err := c.BindJSON(&Movie); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := initializer.CreateMovie(&Movie); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create Movie"})
		return
	}

	c.JSON(201, gin.H{"message": "Movie created successfully"})
}

func GetMovies(c *gin.Context) {
	cursor, err := initializer.TestCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve Movies"})
		return
	}
	defer cursor.Close(context.Background())

	var Movies []models.Movie

	if err := cursor.All(context.Background(), &Movies); err != nil {
		log.Println("Error decoding Movies:", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve Movies"})
		return
	}

	c.JSON(200, Movies)
}
func UpdateMovie(c *gin.Context) {
	title := c.Param("title")

	var body models.Movie
	c.BindJSON(&body)

	filter := bson.M{"Title": title}

	update := bson.M{"$set": bson.M{"Title": body.Title, "Year": body.Year}}

	result, err := initializer.TestCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error decoding Movie:", err)
		return
	}

	c.JSON(200, gin.H{"message": "Movie Updated successfully!", "Movie": result})

}

func DeleteMovie(c *gin.Context) {
	title := c.Param("title")

	result, err := initializer.TestCollection.DeleteMany(context.Background(), bson.M{"Title": title})
	log.Println("Delete result:", result)

	if err != nil {
		log.Println("Error deleting Movie:", err)
		return
	}

	if result.DeletedCount == 0 {
		log.Println("Delete Count:", result.DeletedCount)

		c.JSON(400, gin.H{"error": "Movie Not Found"})
		return
	}
	c.JSON(200, gin.H{"message": "Movie deleted successfully"})

}
