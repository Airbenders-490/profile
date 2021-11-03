package app

import (
	schoolHttp "github.com/airbenders/profile/School/delivery/http"
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/gin-gonic/gin"
)

func mapStudentURLs(h *http.StudentHandler, r *gin.Engine) {
	r.GET("/student/:id", h.GetByID)
	r.POST("/student", h.Create)
	r.PUT("/student/:id", h.Update)
	r.DELETE("/student/:id", h.Delete)
}

// {"AI123", "Adam"}
// {"BI123", "Adam"}

func mapSchoolURLs(h *schoolHttp.SchoolHandler, r *gin.Engine) {
	r.GET("/school", h.SearchStudentSchool)
	r.POST("/school/confirm", h.SendConfirmationMail)
	r.GET("/school/confirmation", h.ConfirmSchoolRegistration)
}