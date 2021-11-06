package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" || authToken[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
			c.Abort()
			return
		}
		client := &http.Client{}
		url := "http://localhost:3000/api/validate"
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "some error occurred while making a request",
			})
			log.Println(err.Error())
			c.Abort()
			return
		}

		request.Header.Set("Authorization", authToken)
		response, err := client.Do(request)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "some error occurred while confirming the token",
			})
			log.Println(err.Error())
			c.Abort()
			return
		}
		if response.StatusCode == 200 {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
			c.Abort()
			return
		}
	}
}
