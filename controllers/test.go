package controllers

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/karthikeyan-amk/go-curd/initializer"
    "github.com/karthikeyan-amk/go-curd/models"
    "github.com/redis/go-redis/v9" // Import Redis package
    "go.mongodb.org/mongo-driver/bson"
)

var redisClient *redis.Client

// InitializeRedisClient initializes Redis client
func InitializeRedisClient() {
    redisClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis server address
        Password: "",               // No password set
        DB:       0,                // Default DB
    })
		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
				panic(err)
		}
}

// SetMovieInCache sets movie data in Redis cache
func SetMovieInCache(movie models.Movie) error {
    ctx := context.Background()
    movieJSON, err := bson.Marshal(movie)
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
    if err := bson.Unmarshal(val, &movie); err != nil {
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

	// GetMovies retrieves all movies from MongoDB
func GetMovies(c *gin.Context) {
	// Find all movies in the MongoDB collection
	cursor, err := initializer.TestCollection.Find(context.Background(), bson.M{})
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies"})
			return
	}
	defer cursor.Close(context.Background())

	var movies []models.Movie

	// Iterate over the cursor and decode each document into a Movie struct
	for cursor.Next(context.Background()) {
			var movie models.Movie
			if err := cursor.Decode(&movie); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movie"})
					return
			}
			movies = append(movies, movie)
	}

	// Check if there was any error during cursor iteration
	if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
			return
	}

	// Return the list of movies
	c.JSON(http.StatusOK, movies)
}

// GetMovie retrieves a single movie by its title from MongoDB or Redis cache
func GetMovie(c *gin.Context) {
	title := c.Query("title") // Retrieve title parameter from query

	// Try to get the movie from Redis cache
	movieFromCache, err := GetMovieFromCache(title)
	if err == nil { // If movie found in cache
			c.JSON(http.StatusOK, movieFromCache)
			return
	}

	// If movie not found in cache, fetch from MongoDB
	var movieFromDB models.Movie
	filter := bson.M{"title": title}
	if err := initializer.TestCollection.FindOne(context.Background(), filter).Decode(&movieFromDB); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in mondodb"})
			return
	}

	// Set movie in Redis cache
	if err := SetMovieInCache(movieFromDB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set movie in cache"})
			return
	}

	// Return the movie from MongoDB
	c.JSON(http.StatusOK, movieFromDB)
}


// UpdateMovie updates a movie in MongoDB and Redis cache if it exists in both
func UpdateMovie(c *gin.Context) {
	title := c.Query("title")
	var updatedMovie models.Movie
	if err := c.BindJSON(&updatedMovie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
	}

	// Check if the movie exists in Redis cache
	movieFromCache, err := GetMovieFromCache(title)
	if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in Redis cache"})
			return
	}

	// Check if the movie exists in MongoDB
	var movieFromDB models.Movie
	filter := bson.M{"title": title}
	if err := initializer.TestCollection.FindOne(context.Background(), filter).Decode(&movieFromDB); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in MongoDB"})
			return
	}

	// Update the movie in Redis cache
	movieFromCache.Year = updatedMovie.Year
	if err := SetMovieInCache(movieFromCache); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie in Redis cache"})
			return
	}

	// Update the movie in MongoDB
	update := bson.M{"$set": bson.M{"year": updatedMovie.Year}}
	if _, err := initializer.TestCollection.UpdateOne(context.Background(), filter, update); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie in MongoDB"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}

// DeleteMovie deletes a movie from MongoDB and Redis cache if it exists in both
func DeleteMovie(c *gin.Context) {
	title := c.Query("title")

	// Check if the movie exists in Redis cache
	_, err := GetMovieFromCache(title)
	if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in Redis cache"})
			return
	}

	// Check if the movie exists in MongoDB
	var movieFromDB models.Movie
	filter := bson.M{"title": title}
	if err := initializer.TestCollection.FindOne(context.Background(), filter).Decode(&movieFromDB); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in MongoDB"})
			return
	}

	// Delete the movie from Redis cache
	ctx := context.Background()
	if err := redisClient.Del(ctx, title).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie from cache"})
			return
	}

	// Delete the movie from MongoDB
	if _, err := initializer.TestCollection.DeleteOne(context.Background(), filter); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie from MongoDB"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
