package models

type Movie struct {
	Title string `json:"title" bson:"title"`
	Year  int    `json:"year" bson:"year"`
}
