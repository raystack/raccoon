package serialization

import "github.com/stretchr/testify/mock"

type MockSerializer struct {
	mock.Mock
}

func (ms *MockSerializer) Serialize(m interface{}) ([]byte, error) {
	args := ms.Called(m)

	return []byte(args.String(0)), args.Error(1)
}
