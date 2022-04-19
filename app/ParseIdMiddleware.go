package app

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type ClaimsParser interface {
	ParseClaimsMiddleware() gin.HandlerFunc
}

type parseClaims struct{}

// NewParseClaimsMiddleware is a constructor
func NewParseClaimsMiddleware() ClaimsParser {
	return &parseClaims{}
}

// ParseClaimsMiddleware add the id to context
func (h *parseClaims) ParseClaimsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.RegisteredClaims
		c.Set("loggedID", claims.Subject)
		c.Next()
	}
}
