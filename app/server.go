package app

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"

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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Server is a constructor. Returns the router after mapping all the urls
func Server(
	studentHandler *http.StudentHandler,
	schoolHandler *http2.SchoolHandler,
	tagHandler *http3.TagHandler,
	reviewHandler *http4.ReviewHandler,
	mw Middleware) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	mapStudentURLs(mw, studentHandler, router)
	mapSchoolURLs(mw, schoolHandler, router)
	mapTagURLs(tagHandler, router)
	mapReviewURLs(mw, reviewHandler, router)
	return router
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Start runs the server
// todo: refactor and breakdown
func Start() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(os.Getenv("DATABASE_URL"))
		log.Fatalln("db failed", err)
	}
	conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	failOnError(err, "can't connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"profile", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "can't create exchange")
	mm := usecase.NewMessagingManager(ch)
	studentRepository := repository.NewStudentRepository(pool)
	reviewRepository := repository4.NewReviewRepository(pool)
	studentUseCase := usecase.NewStudentUseCase(mm, studentRepository, reviewRepository, time.Second)
	studentHandler := http.NewStudentHandler(studentUseCase)
	go studentUseCase.CreateStudentTopic()
	go studentUseCase.UpdateStudentTopic()
	go studentUseCase.DeleteStudentTopic()
	schoolRepository := repository2.NewSchoolRepository(pool)
	mail := utils.NewSimpleMail()
	schoolUseCase := usecase2.NewSchoolUseCase(schoolRepository, studentRepository, mail, time.Second)
	schoolHandler := http2.NewSchoolHandler(schoolUseCase)

	tagRepository := repository3.NewTagRepository(pool)
	tagUseCase := usecase3.NewTagUseCase(tagRepository, time.Second)
	tagHandler := http3.NewTagHandler(tagUseCase)

	reviewUseCase := usecase4.NewReviewUseCase(reviewRepository, studentRepository, time.Second)
	reviewHandler := http4.NewReviewHandler(reviewUseCase)

	mw := NewMiddleware()

	router := Server(studentHandler, schoolHandler, tagHandler, reviewHandler, mw)
	router.Run()
}
