package mocks

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

// StudentUseCase Mock struct
type StudentUseCase struct {
	mock.Mock
}

func (m *StudentUseCase) UpdateStudentTopic() {
	panic("implement me")
}

func (m *StudentUseCase) DeleteStudentTopic() {
	panic("implement me")
}

func (m *StudentUseCase) CreateStudentTopic() {
	panic("implement me")
}

// Create - StudentUseCaseMock
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

// GetByID - StudentUseCaseMock
func (m *StudentUseCase) GetByID(ctx context.Context, id string) (*domain.Student, error) {
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

// Update - StudentUseCaseMock
func (m *StudentUseCase) Update(ctx context.Context, id string, st *domain.Student) (*domain.Student, error) {
	args := m.Called(ctx, id, st)

	var r0 *domain.Student
	if rf, ok := args.Get(0).(func(context.Context, string, *domain.Student) *domain.Student); ok {
		r0 = rf(ctx, id, st)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*domain.Student)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string, *domain.Student) error); ok {
		r1 = rf(ctx, id, st)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

// Delete - StudentUseCaseMock
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

func (m *StudentUseCase) AddClasses(c context.Context, id string, st *domain.Student) error {
	ret := m.Called(c, id, st)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.Student) error); ok {
		r0 = rf(c, id, st)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CompleteClass provides a mock function with given fields: c, id, st
func (m *StudentUseCase) CompleteClass(c context.Context, id string, st *domain.Student) error {
	ret := m.Called(c, id, st)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.Student) error); ok {
		r0 = rf(c, id, st)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *StudentUseCase) RemoveClasses(c context.Context, id string, st *domain.Student) error {
	ret := m.Called(c, id, st)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.Student) error); ok {
		r0 = rf(c, id, st)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (m *StudentUseCase) SearchStudents(ctx context.Context, st *domain.Student) ([]domain.Student, error) {
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
