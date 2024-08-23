package main

import (
	"net/http"

	db "shopping_api/database"
	"shopping_api/pkg/routes"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	db := db.InitDB()

	router := gin.Default()
	router.Use(CORSMiddleware())

	v1 := router.Group("/v1")

	router.GET("/health", func(ctx *gin.Context) {
		response := map[string]string{"message": "Shopping api healthy"}
		ctx.JSON(http.StatusOK, response)
	})

	routes.AddLoginEndpoints(db, v1)
	routes.AddShoppingEndpoints(db, v1)

	router.Run(":8080")
}
