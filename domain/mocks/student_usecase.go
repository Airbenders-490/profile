package mocks

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

type StudentUseCase struct {
	mock.Mock
}

func (m *StudentUseCase) Create(ctx context.Context, st *domain.Student) error {
	args := m.Called(ctx, st)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, *domain.Student) error); ok {
		r0 = rf(ctx, st)
	} else {
		r0 = args.Error(0)
	}
	return r0
}

func (m *StudentUseCase) GetById(ctx context.Context, id string) (*domain.Student, error) {
	args := m.Called(ctx, id)

	var r0 *domain.Student
	if rf, ok := args.Get(0).(func(context.Context, string) *domain.Student); ok {
		r0 = rf(ctx, id)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*domain.Student)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (m *StudentUseCase) Update(ctx context.Context, st *domain.Student) error {
	args := m.Called(ctx, st)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, *domain.Student) error); ok {
		r0 = rf(ctx, st)
	} else {
		r0 = args.Error(0)
	}

	return r0
}

func (m *StudentUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = args.Error(0)
	}
	return r0
}
