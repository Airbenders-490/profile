package app

import (
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/gin-gonic/gin"
)

func mapStudentURLs(h *http.StudentHandler, r *gin.Engine) {
	r.GET("/student/:id", h.GetByID)
	r.POST("/student", h.Create)
	r.PUT("/student/:id", h.Update)
	r.DELETE("/student/:id", h.Delete)
}
