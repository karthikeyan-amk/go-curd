package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/karthikeyan-amk/go-curd/initializer"
	"github.com/karthikeyan-amk/go-curd/models"
	"github.com/redis/go-redis/v9" // Import Redis package
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var redisClient *redis.Client

// InitializeRedisClient initializes Redis client
func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Default DB
	})
}

// SetMovieInCache sets movie data in Redis cache
func SetMovieInCache(movie models.Movie) error {
	ctx := context.Background()
	movieJSON, err := json.Marshal(movie)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, movie.Title, movieJSON, 0).Err()
}

// GetMovieFromCache gets movie data from Redis cache
func GetMovieFromCache(title string) (models.Movie, error) {
	ctx := context.Background()
	val, err := redisClient.Get(ctx, title).Bytes()
	if err != nil {
		return models.Movie{}, err
	}
	var movie models.Movie
	if err := json.Unmarshal(val, &movie); err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}

// CreateMovie creates a new movie in MongoDB and sets it in Redis cache
func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.BindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := initializer.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}
	// Set movie in Redis cache
	if err := SetMovieInCache(movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set movie in cache"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Movie created successfully"})
}

// GetMovie retrieves a movie by title from MongoDB and Redis cache
func GetMovie(c *gin.Context) {
	title := c.Param("title")
	// Get movie from Redis cache
	movie, err := GetMovieFromCache(title)
	if err != nil {
		// If movie not found in cache, fetch from MongoDB
		var movieFromDB models.Movie
		filter := bson.M{"title": title}
		if err := initializer.TestCollection.FindOne(context.Background(), filter).Decode(&movieFromDB); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		// Set movie in Redis cache
		if err := SetMovieInCache(movieFromDB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set movie in cache"})
			return
		}
		c.JSON(http.StatusOK, movieFromDB)
		return
	}
	c.JSON(http.StatusOK, movie)
}

// UpdateMovie updates a movie in MongoDB and Redis cache
func UpdateMovie(c *gin.Context) {
	title := c.Param("title")
	var updatedMovie models.Movie
	if err := c.BindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	filter := bson.M{"title": title}
	update := bson.M{"$set": bson.M{"title": updatedMovie.Title, "year": updatedMovie.Year}}
	var res *mongo.UpdateResult
	var err error
	if res, err = initializer.TestCollection.UpdateOne(context.Background(), filter, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}
	// Update movie in Redis cache
	if err := SetMovieInCache(updatedMovie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set movie in cache"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully", "modified_count": res.ModifiedCount})
}

// DeleteMovie deletes a movie from MongoDB and Redis cache
func DeleteMovie(c *gin.Context) {
	title := c.Param("title")
	filter := bson.M{"title": title}
	if _, err := initializer.TestCollection.DeleteOne(context.Background(), filter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}
	// Delete movie from Redis cache
	ctx := context.Background()
	if err := redisClient.Del(ctx, title).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie from cache"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
