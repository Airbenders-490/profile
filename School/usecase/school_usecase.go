package usecase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"log"
	"strings"
	"time"
)

type schoolUseCase struct {
	r       domain.SchoolRepository
	timeout time.Duration
}

func NewSchoolUseCase(r domain.SchoolRepository, timeout time.Duration) domain.SchoolUseCase {
	return &schoolUseCase{r, timeout}
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
