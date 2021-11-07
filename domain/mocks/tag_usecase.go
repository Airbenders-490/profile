package mocks

import (
	"context"

	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

// TagUseCase Mock struct
type TagUseCase struct {
	mock.Mock
}

// GetAllTags - TagUseCase
func (m *TagUseCase) GetAllTags(ctx context.Context) ([]domain.Tag, error) {
	args := m.Called(ctx)

	var r0 []domain.Tag
	if rf, ok := args.Get(0).(func(context.Context) []domain.Tag); ok {
		r0 = rf(ctx)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).([]domain.Tag)
		}
	}
	var r1 error
	if rf, ok := args.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1

}

// GetByID - StudentUseCaseMock
