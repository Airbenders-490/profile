package app

import (
	"context"
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/airbenders/profile/Student/repository"
	"github.com/airbenders/profile/Student/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

func Server(h *http.StudentHandler) *gin.Engine {
	router := gin.Default()
	mapStudentURLs(h, router)
	return router
}

// Start runs the server
func Start() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	repository := repository.NewStudentRepository(pool)
	useCase := usecase.NewStudentUseCase(repository, time.Second)
	h := http.NewStudentHandler(useCase)
	router := Server(h)
	router.Run()
}
