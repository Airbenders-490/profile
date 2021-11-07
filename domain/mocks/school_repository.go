package mocks

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

// SchoolRepositoryMock struct
type SchoolRepositoryMock struct {
	mock.Mock
}

// SearchByDomain -- SchoolRepositoryMock
func (m *SchoolRepositoryMock) SearchByDomain(ctx context.Context, name string) ([]domain.School, error) {
	args := m.Called(ctx, name)

	var r0 []domain.School
	if rf, ok := args.Get(0).(func(context.Context, string) []domain.School); ok {
		r0 = rf(ctx, name)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).([]domain.School)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

// SaveConfirmationToken -- SchoolRepositoryMock
func (m *SchoolRepositoryMock) SaveConfirmationToken(ctx context.Context, confirmation *domain.Confirmation) error {
	args := m.Called(ctx, confirmation)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, *domain.Confirmation) error); ok {
		r0 = rf(ctx, confirmation)
	} else {
		r0 = args.Error(0)
	}
	return r0
}

// GetConfirmationByToken -- SchoolRepositoryMock
func (m *SchoolRepositoryMock) GetConfirmationByToken(ctx context.Context, token string) (*domain.Confirmation, error) {
	args := m.Called(ctx, token)

	var r0 *domain.Confirmation
	if rf, ok := args.Get(0).(func(context.Context, string) *domain.Confirmation); ok {
		r0 = rf(ctx, token)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*domain.Confirmation)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1

}

// AddSchoolForStudent -- SchoolRepositoryMock
func (m *SchoolRepositoryMock) AddSchoolForStudent(ctx context.Context, stID string, scID string) error {
	args := m.Called(ctx, stID, scID)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, stID, scID)
	} else {
		r0 = args.Error(0)
	}
	return r0
}
