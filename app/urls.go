package app

import (
	reviewHttp "github.com/airbenders/profile/Review/delivery/http"
	schoolHttp "github.com/airbenders/profile/School/delivery/http"
	studentHttp "github.com/airbenders/profile/Student/delivery/http"
	tagHttp "github.com/airbenders/profile/Tag/delivery/http"
	"github.com/gin-gonic/gin"
)

func mapStudentURLs(m Middleware, h *studentHttp.StudentHandler, router *gin.Engine) {
	authorized := router.Group("/api")
	authorized.Use(m.AuthMiddleware())
	const pathStudentID = "/student/:id"
	authorized.GET(pathStudentID, h.GetByID)
	authorized.POST("/student", h.Create)
	authorized.PUT(pathStudentID, h.Update)
	authorized.DELETE(pathStudentID, h.Delete)
	authorized.PUT("/addClasses/:id", h.AddClasses)
	authorized.PUT("/removeClasses/:id", h.RemoveClasses)
	authorized.PUT("/completeClass/:id", h.CompleteClass)
}

func mapSchoolURLs(m Middleware, h *schoolHttp.SchoolHandler, r *gin.Engine) {
	r.GET("school/confirmation", h.ConfirmSchoolRegistration)
	authorized := r.Group("/api")
	authorized.Use(m.AuthMiddleware())
	authorized.GET("/school", h.SearchStudentSchool)
	authorized.POST("/school/confirm", h.SendConfirmationMail)
}

func mapTagURLs(h *tagHttp.TagHandler, r *gin.Engine) {
	r.GET("/api/all-tags", h.GetAllTags)
}

func mapReviewURLs(m Middleware, h *reviewHttp.ReviewHandler, r *gin.Engine) {
	authorized := r.Group("/api")
	authorized.Use(m.AuthMiddleware())
	authorized.POST("/review/:reviewed", h.AddReview)
	authorized.PUT("/review/:reviewed/update", h.EditReview)
	authorized.GET("/reviews-by/:reviewer", h.GetReviewsBy)
}
