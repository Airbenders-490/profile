package app

import (
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/airbenders/profile/Student/repository"
	"github.com/airbenders/profile/Student/usecase"
	"github.com/gin-gonic/gin"
	"time"
)

func  mapURLs(r *gin.Engine) {
	repository := repository.NewStudentRepository()
	useCase := usecase.NewStudentUseCase(repository, time.Second)
	h := http.NewStudentHandler(useCase)
	r.GET("/student/:id", h.GetByID)
	r.POST("/student", h.Create)
	r.PUT("/student/:id", h.Update)
	r.DELETE("/student/:id", h.Delete)
}
