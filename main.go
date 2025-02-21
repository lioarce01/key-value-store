package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var store sync.Map

// SET handler
func setHandler(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	if key == "" || key == "" {
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
	keys := make([]string, 0)
	store.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// HEALTH handler
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func main() {
	r := gin.Default()

	r.GET("/set", setHandler)
	r.GET("/get", getHandler)
	r.GET("/delete", deleteHandler)
	r.GET("/keys", keysHandler)
	r.GET("/health", healthHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failted to start: %v", err)
	}
}
