package app

import (
	reviewHttp "github.com/airbenders/profile/Review/delivery/http"
	schoolHttp "github.com/airbenders/profile/School/delivery/http"
	studentHttp "github.com/airbenders/profile/Student/delivery/http"
	tagHttp "github.com/airbenders/profile/Tag/delivery/http"
	"github.com/airbenders/profile/app/middlwares"
	"github.com/gin-gonic/gin"
)

func studentURLs(h *studentHttp.StudentHandler, authorized *gin.RouterGroup) {
	const pathStudentID = "/student/:id"
	authorized.GET(pathStudentID, h.GetByID)
	authorized.POST("/student", h.Create)
	authorized.PUT(pathStudentID, h.Update)
	authorized.DELETE(pathStudentID, h.Delete)
	authorized.PUT("/addClasses/:id", h.AddClasses)
	authorized.PUT("/removeClasses/:id", h.RemoveClasses)
	authorized.PUT("/completeClasses/:id", h.CompleteAllClasses)
	authorized.GET("/search/", h.SearchStudents)
}

func mapStudentURLsV0(m middlwares.Middleware, h *studentHttp.StudentHandler, router *gin.Engine) {
	authorized := router.Group("/api")
	authorized.Use(m.AuthMiddleware())
	studentURLs(h, authorized)
}

func mapStudentURLsV1(m middlwares.Middleware, parserMW middlwares.ClaimsParser, h *studentHttp.StudentHandler, router *gin.Engine) {
	authorized := router.Group("/api/v1")
	authorized.Use(m.AuthMiddleware())
	authorized.Use(parserMW.ParseClaimsMiddleware())
	studentURLs(h, authorized)
}

func mapSchoolURLsV0(m middlwares.Middleware, h *schoolHttp.SchoolHandler, r *gin.Engine) {
	r.GET("school/confirmation", h.ConfirmSchoolRegistration)
	authorized := r.Group("/api")
	authorized.Use(m.AuthMiddleware())
	schoolURLs(authorized, h)
}

// we can extract these 2 into api versions for better abstraction and maintainability
func mapSchoolURLsV1(m middlwares.Middleware, parserMW middlwares.ClaimsParser, h *schoolHttp.SchoolHandler, r *gin.Engine) {
	//r.GET("school/confirmation", h.ConfirmSchoolRegistration)
	authorized := r.Group("/api/v1")
	authorized.Use(m.AuthMiddleware())
	authorized.Use(parserMW.ParseClaimsMiddleware())
	schoolURLs(authorized, h)
}

func schoolURLs(authorized *gin.RouterGroup, h *schoolHttp.SchoolHandler) {
	authorized.GET("/school", h.SearchStudentSchool)
	authorized.GET("/school/confirm", h.SendConfirmationMail)
}

func mapTagURLs(h *tagHttp.TagHandler, r *gin.Engine) {
	r.GET("/api/all-tags", h.GetAllTags)
}

func reviewURLs(h *reviewHttp.ReviewHandler, authorized *gin.RouterGroup) {
	authorized.POST("/review/:reviewed", h.AddReview)
	authorized.PUT("/review/:reviewed/update", h.EditReview)
	authorized.GET("/reviews-by/:reviewer", h.GetReviewsBy)
}

func mapReviewURLsV0(m middlwares.Middleware, h *reviewHttp.ReviewHandler, r *gin.Engine) {
	authorized := r.Group("/api")
	authorized.Use(m.AuthMiddleware())
	reviewURLs(h, authorized)
}

func mapReviewURLsV1(m middlwares.Middleware, parserMW middlwares.ClaimsParser, h *reviewHttp.ReviewHandler, r *gin.Engine) {
	authorized := r.Group("/api/v1")
	authorized.Use(m.AuthMiddleware())
	authorized.Use(parserMW.ParseClaimsMiddleware())
	reviewURLs(h, authorized)
}
