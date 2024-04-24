package main

import (
	"github.com/gin-gonic/gin"
	"github.com/karthikeyan-amk/go-curd/controllers"
	"github.com/karthikeyan-amk/go-curd/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.Database()
}
func main() {
	r := gin.Default()
	r.GET("/getmovie/:title", controllers.GetMovie)
	r.POST("/createmovie", controllers.CreateMovie)
	r.PUT("/updatemovie/:title", controllers.UpdateMovie)
	r.DELETE("/deletemovie/:title", controllers.DeleteMovie)

	r.Run(":3000")
}
