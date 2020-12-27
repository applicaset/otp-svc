package otpsvc_test

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockSMSService struct {
	mock.Mock
	LastMessage string
}

func (svc *MockSMSService) SendMessage(ctx context.Context, phoneNumber, message string) error {
	svc.LastMessage = message

	args := svc.Called(ctx, phoneNumber, message)

	return args.Error(0)
}
