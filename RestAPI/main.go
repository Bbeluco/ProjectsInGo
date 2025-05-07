package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type album struct {
	ID 		string  `json:"id"`
	Title	string  `json:"title"`
	Artist	string  `json:"artist"`
	Price	float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Barril de RAP vol. 2", Artist: "Barril de rap", Price: 100.99},
	{ID: "2", Title: "Slim shady LP", Artist: "Eminem", Price: 23.54},
	{ID: "3", Title: "The Eminem Show", Artist: "Eminem", Price: 96.59},
}


func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		log.Error("Unable to parse JSON", err.Error());
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while parsing"})
		return;
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return;
		}
	}

	// c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"});
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while parsing"})
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("albums/:id", getAlbumByID)
	router.POST("albums", postAlbums)

	router.Run("localhost:8080")
}