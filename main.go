package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

var store sync.Map

// SET handler
func setHandler(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		return
	}

	store.Store(key, value)
	c.JSON(http.StatusOK, gin.H{"message": "key set successfully"})
}

// GET handler
func getHandler(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	value, ok := store.Load(key)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}

// DELETE handler
func deleteHandler(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	_, ok := store.Load(key)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}

	store.Delete(key)
	c.JSON(http.StatusOK, gin.H{"message": "key deleted successfully"})
}

// KEYS handler
func keysHandler(c *gin.Context) {
	var keys []string
	store.Range(func(k, v interface{}) bool {
		keys = append(keys, k.(string))
		return true
	})
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// HEALTH handler
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// RATE LIMITER
var limiter = rate.NewLimiter(1, 5)

func rateMiddleware(c *gin.Context) {
	if limiter.Allow() == false {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	r := gin.Default()
	r.Use(rateMiddleware)

	r.GET("/set", setHandler)
	r.GET("/get", getHandler)
	r.GET("/delete", deleteHandler)
	r.GET("/keys", keysHandler)
	r.GET("/health", healthHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failted to start: %v", err)
	}
}
