package otpsvc_test

import (
	"context"
	otpsvc "github.com/applicaset/otp-svc"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
	LastItem otpsvc.Entity
}

func (repo *MockRepository) Create(ctx context.Context, entity otpsvc.Entity) error {
	repo.LastItem = entity

	args := repo.Called(ctx, entity)

	return args.Error(0)
}

func (repo *MockRepository) Find(ctx context.Context, otpUUID string) (*otpsvc.Entity, error) {
	args := repo.Called(ctx, otpUUID)

	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}

	return res.(*otpsvc.Entity), args.Error(1)
}
