package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MiddlewareMock struct {
	mock.Mock
}

func (m *MiddlewareMock) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("id")
		c.Set("loggedID", id)
		c.Next()
	}
}
