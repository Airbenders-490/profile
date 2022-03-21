package usecase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils"
	"github.com/airbenders/profile/utils/errors"
	"github.com/google/uuid"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"
)

type schoolUseCase struct {
	r       domain.SchoolRepository
	str     domain.StudentRepository
	mailer  utils.Mailer
	timeout time.Duration
}

// NewSchoolUseCase is the constructor
func NewSchoolUseCase(r domain.SchoolRepository, str domain.StudentRepository, mailer utils.Mailer, timeout time.Duration) domain.SchoolUseCase {
	return &schoolUseCase{r, str, mailer, timeout}
}

// SearchSchoolByDomain doesn't do much processing. just returns the school slice if received from the repository.
// returns nil and 404 error if there are no schools
func (s *schoolUseCase) SearchSchoolByDomain(c context.Context, domainName string) ([]domain.School, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	domainNameLike := parseDomainName(domainName)
	schools, err := s.r.SearchByDomain(ctx, domainNameLike)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	if len(schools) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no school with domain name %s exists", domainName))
	}

	return schools, nil
}

func parseDomainName(name string) string {
	var domains []string
	// first add the same domain as we got: e.g. live.concordia.ca to it
	domains = append(domains, name)
	count := strings.Count(name, ".")
	// live.concordia.ca has 2
	if count > 1 {
		if first := strings.Index(name, "."); first != len(name) {
			// extract the latter half like concordia.ca and add it to domains
			name = name[first+1:]
			domains = append(domains, name)
		}
	}
	// string builder
	var domainNameLike bytes.Buffer
	// write the first one. This shouldn't be an error
	_, err := domainNameLike.WriteString(domains[0])
	if err != nil {
		log.Println("ISSUE!")
	}
	// for the rest, add them with an 'or' operator
	for _, domainName := range domains[1:] {
		domainNameLike.WriteString(fmt.Sprintf("|%s", domainName))
	}
	return domainNameLike.String()
}

// SendConfirmation sends an email to the student's school email address with a generated token
// also stores the token in the repository with the student's and school's IDs for confirmation later
func (s *schoolUseCase) SendConfirmation(c context.Context, st *domain.Student, email string, school *domain.School) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	student, err := s.str.GetByID(ctx, st.ID)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(student, &domain.Student{}) {
		return errors.NewNotFoundError("student not found")
	}
	if student.School != nil {
		return errors.NewBadRequestError("school already confirmed")
	}

	token := uuid.New().String()
	domainName := os.Getenv("DOMAIN")
	if domainName == "" {
		log.Fatalln("Domain name not provided")
	}
	confirmationURL := fmt.Sprintf("%s/school/confirmation", domainName)
	url := fmt.Sprintf("%s?token=%s", confirmationURL, token)

	confirmation := &domain.Confirmation{
		Token:     token,
		School:    *school,
		Student:   *st,
		CreatedAt: time.Now(),
	}

	err = s.r.SaveConfirmationToken(ctx, confirmation)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	body := createEmailBody(student.FirstName, school.Name, url)
	return s.mailer.SendSimpleMail(email, body)
}

func createEmailBody(name, school, url string) []byte {
	t, err := template.ParseFiles("static/confirmation_template.html")
	if err != nil {
		t, err = template.ParseFiles("../../static/confirmation_template.html")
		if err != nil {
			log.Fatal("no email confirmation template")
		}
	}
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Smarties app confirmation\n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Name   string
		School string
		Email  string
	}{
		Name:   name,
		School: school,
		Email:  url,
	})

	return body.Bytes()
}

// ConfirmSchoolEnrollment checks if the record for the token exists in the repository.
// if it does, it checks to ensure it's more than 24 hours old
func (s *schoolUseCase) ConfirmSchoolEnrollment(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	confirmation, err := s.r.GetConfirmationByToken(ctx, token)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	if reflect.DeepEqual(*confirmation, domain.Confirmation{}) {
		return errors.NewNotFoundError("invalid token")
	}
	if confirmation.CreatedAt.Add(time.Hour * 24).Before(time.Now()) {
		return errors.NewBadRequestError("token already expired")
	}

	err = s.r.AddSchoolForStudent(ctx, confirmation.Student.ID, confirmation.School.ID)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}
