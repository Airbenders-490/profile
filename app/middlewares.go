package app

import (
	"github.com/airbenders/profile/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// Middleware defines the contracts
type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type middleware struct {}

// NewMiddleware is a constructor
func NewMiddleware() Middleware {
	return &middleware{}
}

// AuthMiddleware checks if it has a jwt token and then requests the auth service to verify it for us
func (h *middleware) AuthMiddleware() gin.HandlerFunc {
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
			jwtToken := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1)

			token, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})
			if err != nil {
				c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid token. I think"))
				return
			}
			var loggedID string
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				iss := claims["iss"]
				if f, ok := iss.(string); ok {
					loggedID = f
				}
			} else {
				c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid token. I think"))
				return
			}
			c.Set("loggedID", loggedID)
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
