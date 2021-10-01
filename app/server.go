package app

import (
	"github.com/gin-gonic/gin"
)

func server() *gin.Engine {
	router := gin.Default()
	return router
}

func Start() {
	router := server()
	router.Run()
}
