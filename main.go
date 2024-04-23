package main

import (
	"github.com/gin-gonic/gin"
	"github.com/karthikeyan-amk/go-curd/controllers"
)

func main() {
	r := gin.Default()
	r.GET("/getmovies", controllers.GetMovies)
	r.POST("/createmovie", controllers.CreateMovie)
	r.PUT("/updatemovie/:id", controllers.UpdateMovie)
	r.DELETE("deletemovie/:id", controllers.DeleteMovie)

	r.Run(":3000")
}
