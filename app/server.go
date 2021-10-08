package app

import (
	"github.com/gin-gonic/gin"
)

func server() *gin.Engine {
	router := gin.Default()
	mapURLs(router)
	return router
}

// Start runs the server
func Start() {
	router := server()
	router.Run()
}
