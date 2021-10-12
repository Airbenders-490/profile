package app

import (
	"context"
	http2 "github.com/airbenders/profile/School/delivery/http"
	repository2 "github.com/airbenders/profile/School/repository"
	usecase2 "github.com/airbenders/profile/School/usecase"
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/airbenders/profile/Student/repository"
	"github.com/airbenders/profile/Student/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

func Server(stundetHandler *http.StudentHandler, schoolHandler *http2.SchoolHandler) *gin.Engine {
	router := gin.Default()
	mapStudentURLs(stundetHandler, router)
	mapSchoolURLs(schoolHandler, router)
	return router
}

// Start runs the server
func Start() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	studentRepository := repository.NewStudentRepository(pool)
	studentUseCase := usecase.NewStudentUseCase(studentRepository, time.Second)
	studentHandler := http.NewStudentHandler(studentUseCase)

	schoolRepository := repository2.NewSchoolRepository(pool)
	schoolUseCase := usecase2.NewSchoolUseCase(schoolRepository, time.Second)
	schoolHandler := http2.NewSchoolHandler(schoolUseCase)

	router := Server(studentHandler, schoolHandler)
	router.Run()
}
