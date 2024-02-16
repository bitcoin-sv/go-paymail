package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// index basic request to /
// nolint: revive // do not check for unused param required by interface
func index(c *gin.Context) {
	responseData := map[string]interface{}{"message": "Welcome to the Paymail Server ✌(◕‿-)✌"}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, responseData)
}

// health is a basic request to return a health response
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

// notFound handles all 404 requests
// nolint: revive // do not check for unused param required by interface
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "request not found")
}

// methodNotAllowed handles all 405 requests
func methodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, "method"+c.Request.Method+" not allowed")
}
