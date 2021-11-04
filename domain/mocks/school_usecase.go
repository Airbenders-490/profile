package mocks

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

type SchoolUseCase struct {
	mock.Mock
}

func (m *SchoolUseCase) SearchSchoolByDomain(c context.Context, domainName string) ([]domain.School, error){
	args := m.Called(c, domainName)

	var r0 []domain.School
	if rf, ok := args.Get(0).(func(context.Context, string) []domain.School); ok{
		r0 = rf(c, domainName)
	}else{
		if args.Get(0) != nil{
			r0 = args.Get(0).([]domain.School)
		}
	}
	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, domainName)
	} else {
		r1 = args.Error(1)
	}
	return r0, r1
}

func (m *SchoolUseCase) SendConfirmation(c context.Context, st *domain.Student, email string, school *domain.School) error {
	args := m.Called(c, st, email, school)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, *domain.Student, string, *domain.School) error); ok{
		r0 = rf(c, st, email, school)
	} else {
		r0 = args.Error(0)
	}
	return r0
}

func (m *SchoolUseCase) ConfirmSchoolEnrollment(c context.Context, token string) error {
	args := m.Called(c, token)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, string) error); ok{
		r0 = rf(c, token)
	} else {
		r0 = args.Error(0)
	}

	return r0
}