package otp_svc

import "context"

type Repository interface {
	Create(ctx context.Context, entity Entity) (err error)
	Find(ctx context.Context, otpUUID string) (res *Entity, err error)
}
