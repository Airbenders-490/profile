package mocks

import (
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

type SchoolUseCase struct {
	mock.Mock
}

func (s *SchoolUseCase) SearchSchoolByDomain(c context.Context, domainName string) ([]domain.School, error) {
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

