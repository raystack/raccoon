package collection

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCollector struct {
	mock.Mock
}

func (m *MockCollector) Collect(ctx context.Context, req *CollectRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
