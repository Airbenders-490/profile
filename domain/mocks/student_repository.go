package mocks

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

// StudentRepositoryMock struct
type StudentRepositoryMock struct {
	mock.Mock
}

// Create -- StudentRepositoryMock
func (m *StudentRepositoryMock) Create(ctx context.Context, id string, st *domain.Student) error {
	args := m.Called(ctx, id, st)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, string, *domain.Student) error); ok {
		r0 = rf(ctx, id, st)
	} else {
		r0 = args.Error(0)
	}
	return r0
}

// GetByID -- StudentRepositoryMock
func (m *StudentRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Student, error) {
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

// Update -- StudentRepositoryMock
func (m *StudentRepositoryMock) Update(ctx context.Context, st *domain.Student) error {
	args := m.Called(ctx, st)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, *domain.Student) error); ok {
		r0 = rf(ctx, st)
	} else {
		r0 = args.Error(0)
	}

	return r0
}

// Delete -- StudentRepositoryMock
func (m *StudentRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	var r0 error
	if rf, ok := args.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = args.Error(0)
	}
	return r0
}

func (m *StudentRepositoryMock) UpdateClasses(c context.Context, st *domain.Student) error {
	ret := m.Called(c, st)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Student) error); ok {
		r0 = rf(c, st)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (m *StudentRepositoryMock) SearchStudents(ctx context.Context, st *domain.Student) ([]domain.Student, error) {
	ret := m.Called(ctx, st)

	var r0 []domain.Student
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Student) []domain.Student); ok {
		r0 = rf(ctx, st)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Student)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.Student) error); ok {
		r1 = rf(ctx, st)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}