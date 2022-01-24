package mocks

import (
	"github.com/stretchr/testify/mock"
)

// SimpleMail mock struct
type SimpleMail struct {
	mock.Mock
}

// SendSimpleMail mock function
func (m SimpleMail) SendSimpleMail(to string, body []byte) error {
	args := m.Called(to, body)

	var r0 error
	if rf, ok := args.Get(0).(func(string, []byte) error); ok {
		r0 = rf(to, body)
	} else {
		r0 = args.Error(0)
	}
	return r0
}
