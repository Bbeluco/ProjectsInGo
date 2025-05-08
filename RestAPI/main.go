package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Barril de RAP vol. 2", Artist: "Barril de rap", Price: 100.99},
	{ID: "2", Title: "Slim shady LP", Artist: "Eminem", Price: 23.54},
	{ID: "3", Title: "The Eminem Show", Artist: "Eminem", Price: 96.59},
}

var limiter = rate.NewLimiter(1, 5)

var secretKey = "my-secret-key"

func rateLimiter(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		log.Error("Unable to parse JSON", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while parsing"})
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func handleLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var role string
	if username == "admin" && password == "password" {
		role = "admin"
	} else if username == "user" && password == "password" {
		role = "user"
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(jwt.MapClaims)
		role := claims["role"].(string)

		if role != "admin" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func specialAdminContent(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "You're an admin!!",
	})
}

func main() {
	router := gin.Default()

	router.Use(rateLimiter)
	router.POST("/login", handleLogin)
	router.GET("/albums", getAlbums)
	router.GET("albums/:id", getAlbumByID)
	router.Use(authMiddleware())

	router.POST("albums", postAlbums)
	router.GET("special", adminMiddleware(), specialAdminContent)

	router.Run("localhost:8080")
}
