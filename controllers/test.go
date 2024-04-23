package controllers

import (
	"context"
	"log"
	"net/http"

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
	// Get movie title from URL
	title := c.Param("title")

	// Bind JSON data to Movie struct
	var updatedMovie models.Movie
	if err := c.BindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Update movie in the database based on the movie title
	filter := bson.M{"title": title}
	update := bson.M{"$set": bson.M{"year": updatedMovie.Year}}

	_, err := initializer.TestCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}

func DeleteMovie(c *gin.Context) {
	// Get movie title from URL
	title := c.Param("title")

	// Delete movie from the database based on the movie title
	filter := bson.M{"title": title}

	_, err := initializer.TestCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}