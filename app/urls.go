package app

import (
	reviewHttp "github.com/airbenders/profile/Review/delivery/http"
	schoolHttp "github.com/airbenders/profile/School/delivery/http"
	studentHttp "github.com/airbenders/profile/Student/delivery/http"
	tagHttp "github.com/airbenders/profile/Tag/delivery/http"
	"github.com/gin-gonic/gin"
)

func mapStudentURLs(h *studentHttp.StudentHandler, r *gin.Engine) {
	r.GET("/student/:id", h.GetByID)
	r.POST("/student", h.Create)
	r.PUT("/student/:id", h.Update)
	r.DELETE("/student/:id", h.Delete)
}

func mapSchoolURLs(h *schoolHttp.SchoolHandler, r *gin.Engine) {
	r.GET("/school", h.SearchStudentSchool)
	r.POST("/school/confirm", h.SendConfirmationMail)
	r.GET("/school/confirmation", h.ConfirmSchoolRegistration)
}

func mapTagURLs(h *tagHttp.TagHandler, r *gin.Engine) {
	r.GET("/all-tags", h.GetAllTags)
}

func mapReviewURLs(h *reviewHttp.ReviewHandler, r *gin.Engine) {
	r.POST("/review/:reviewed", h.AddReview)
	r.PUT("/review/:reviewed/update", h.EditReview)
	r.GET("/reviews-by/:reviewer", h.GetReviewsBy)
}
