package main

import (
	"github.com/gin-gonic/gin"
	"github.com/karthikeyan-amk/go-curd/controllers"
	"github.com/karthikeyan-amk/go-curd/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.Database()
	controllers.InitializeRedisClient()
}
func main() {
	r := gin.Default()
	r.GET("/getmovie", controllers.GetMovie)
	r.GET("/getmovies", controllers.GetMovies)
	r.POST("/createmovie", controllers.CreateMovie)
	r.PUT("/updatemovie", controllers.UpdateMovie)
	r.DELETE("/deletemovie", controllers.DeleteMovie)
	r.Run(":3000")
}
