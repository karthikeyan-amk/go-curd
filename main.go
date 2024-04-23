package main

import (
	"github.com/gin-gonic/gin"
	"github.com/karthikeyan-amk/go-curd/controllers"
)

func main() {
	r := gin.Default()
	r.GET("/getmovies", controllers.GetMovies)
	r.POST("/createmovie", controllers.CreateMovie)
	r.PUT("/updatemovie/:title", controllers.UpdateMovie)
	r.DELETE("/deletemovie/:title", controllers.DeleteMovie)

	r.Run(":3000")
}
