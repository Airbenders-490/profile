package app

import (
	"context"
	http4 "github.com/airbenders/profile/Review/delivery/http"
	repository4 "github.com/airbenders/profile/Review/repository"
	usecase4 "github.com/airbenders/profile/Review/usecase"
	http2 "github.com/airbenders/profile/School/delivery/http"
	repository2 "github.com/airbenders/profile/School/repository"
	usecase2 "github.com/airbenders/profile/School/usecase"
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/airbenders/profile/Student/repository"
	"github.com/airbenders/profile/Student/usecase"
	http3 "github.com/airbenders/profile/Tag/delivery/http"
	repository3 "github.com/airbenders/profile/Tag/repository"
	usecase3 "github.com/airbenders/profile/Tag/usecase"
	"github.com/airbenders/profile/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

// Server is a constructor. Returns the router after mapping all the urls
func Server(
	studentHandler *http.StudentHandler,
	schoolHandler *http2.SchoolHandler,
	tagHandler *http3.TagHandler,
	reviewHandler *http4.ReviewHandler) *gin.Engine {
	router := gin.Default()
	mapStudentURLs(studentHandler, router)
	mapSchoolURLs(schoolHandler, router)
	mapTagURLs(tagHandler, router)
	mapReviewURLs(reviewHandler, router)
	return router
}

// Start runs the server
func Start() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	studentRepository := repository.NewStudentRepository(pool)
	reviewRepository := repository4.NewReviewRepository(pool)
	studentUseCase := usecase.NewStudentUseCase(studentRepository, reviewRepository, time.Second)
	studentHandler := http.NewStudentHandler(studentUseCase)

	schoolRepository := repository2.NewSchoolRepository(pool)
	mail := utils.NewSimpleMail()
	schoolUseCase := usecase2.NewSchoolUseCase(schoolRepository, studentRepository, mail, time.Second)
	schoolHandler := http2.NewSchoolHandler(schoolUseCase)

	tagRepository := repository3.NewTagRepository(pool)
	tagUseCase := usecase3.NewTagUseCase(tagRepository, time.Second)
	tagHandler := http3.NewTagHandler(tagUseCase)

	reviewUseCase := usecase4.NewReviewUseCase(reviewRepository, studentRepository, time.Second)
	reviewHandler := http4.NewReviewHandler(reviewUseCase)

	router := Server(studentHandler, schoolHandler, tagHandler, reviewHandler)
	router.Run()
}
