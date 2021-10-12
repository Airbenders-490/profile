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
	"strings"
	"text/template"
	"time"
)

type schoolUseCase struct {
	r       domain.SchoolRepository
	str domain.StudentRepository
	mailer utils.Mailer
	timeout time.Duration
}

func NewSchoolUseCase(r domain.SchoolRepository, str domain.StudentRepository, mailer utils.Mailer, timeout time.Duration) domain.SchoolUseCase {
	return &schoolUseCase{r, str, mailer, timeout}
}

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

func (s *schoolUseCase) SendConfirmation(c context.Context, st *domain.Student, email string, school *domain.School) error {
	ctx, cancel :=  context.WithTimeout(c, s.timeout)
	defer cancel()

	token := uuid.New().String()
	domainName := os.Getenv("DOMAIN")
	if domainName == "" {
		log.Fatalln("Domain name not provided")
	}
	confirmationUrl := fmt.Sprintf("%s/school/confirmation", domainName)
	url := fmt.Sprintf("%s?token=%s", confirmationUrl, token)

	confirmation := &domain.Confirmation{
		Token:     token,
		School:    *school,
		Student:   *st,
		CreatedAt: time.Now(),
	}

	err := s.r.SaveConfirmationToken(ctx, confirmation)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	body := createEmailBody(st.FirstName, school.Name, url)
	return s.mailer.SendSimpleMail(email, body)
}

func createEmailBody(name, school, url string) []byte {
	t, err := template.ParseFiles("static/confirmation_template.html")
	if err != nil {
		log.Fatalln(err)
	}
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	fmt.Println(url)
	t.Execute(&body, struct {
		Name string
		School string
		Email string
	}{
		Name: name,
		School: school,
		Email: url,
	})

	return body.Bytes()
}
