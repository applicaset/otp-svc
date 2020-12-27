package otpsvc_test

import (
	"context"
	otpsvc "github.com/applicaset/otp-svc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"time"
)

func TestService_SendOTP(t *testing.T) {
	repo := new(MockRepository)
	smsSvc := new(MockSMSService)
	ctx := context.Background()
	phoneNumber := "+1234567890"

	t.Run("Default", func(t *testing.T) {
		svc := otpsvc.New(repo, smsSvc)

		repo.On("Create", ctx, mock.Anything).Return(nil)

		smsSvc.On("SendMessage", ctx, phoneNumber, mock.Anything).Return(nil)

		res, err := svc.SendOTP(ctx, otpsvc.SendOTPRequest{PhoneNumber: phoneNumber})
		assert.NoError(t, err)
		assert.NotNil(t, res)

		repo.On("Find", ctx, res.OTPUUID).Return(&repo.LastItem, nil)

		res2, err := svc.VerifyOTP(ctx, otpsvc.VerifyOTPRequest{
			OTPUUID: res.OTPUUID,
			PinCode: smsSvc.LastMessage,
		})

		assert.NoError(t, err)
		assert.Equal(t, phoneNumber, res2.PhoneNumber)
	})

	t.Run("Invalid Pin Code", func(t *testing.T) {
		svc := otpsvc.New(repo, smsSvc)

		repo.On("Create", ctx, mock.Anything).Return(nil)

		smsSvc.On("SendMessage", ctx, phoneNumber, mock.Anything).Return(nil)

		res, err := svc.SendOTP(ctx, otpsvc.SendOTPRequest{PhoneNumber: phoneNumber})
		assert.NoError(t, err)
		assert.NotNil(t, res)

		repo.On("Find", ctx, res.OTPUUID).Return(&repo.LastItem, nil)

		res2, err := svc.VerifyOTP(ctx, otpsvc.VerifyOTPRequest{
			OTPUUID: res.OTPUUID,
			PinCode: "Invalid",
		})

		var terr otpsvc.ErrInvalidPinCode
		assert.True(t, errors.As(err, &terr))
		assert.Nil(t, res2)
	})

	t.Run("Send Expired OTP", func(t *testing.T) {
		svc := otpsvc.New(repo, smsSvc)

		repo.On("Create", ctx, mock.Anything).Return(nil)

		smsSvc.On("SendMessage", ctx, phoneNumber, mock.Anything).Return(nil)

		res, err := svc.SendOTP(ctx, otpsvc.SendOTPRequest{PhoneNumber: phoneNumber})
		assert.NoError(t, err)
		assert.NotNil(t, res)

		entity := repo.LastItem
		entity.ExpiresAt = time.Now()

		repo.On("Find", ctx, res.OTPUUID).Return(&entity, nil)

		res2, err := svc.VerifyOTP(ctx, otpsvc.VerifyOTPRequest{
			OTPUUID: res.OTPUUID,
			PinCode: "Invalid",
		})

		var terr otpsvc.ErrOTPNotFoundOrExpired
		assert.True(t, errors.As(err, &terr))
		assert.Nil(t, res2)
	})

	t.Run("WithGoogleAutomaticSMSVerificationTemplate", func(t *testing.T) {
		svc := otpsvc.New(repo, smsSvc, otpsvc.WithGoogleAutomaticSMSVerificationTemplate())

		repo.On("Create", ctx, mock.Anything).Return(nil)

		smsSvc.On("SendMessage", ctx, phoneNumber, mock.Anything).Return(nil)

		hash := "FA+9qCX9VSu"

		res, err := svc.SendOTP(
			ctx,
			otpsvc.SendOTPRequest{PhoneNumber: phoneNumber},
			otpsvc.WithGoogleAutomaticSMSVerification(hash),
		)

		assert.NoError(t, err)
		assert.NotNil(t, res)

		repo.On("Find", ctx, res.OTPUUID).Return(&repo.LastItem, nil)

		msg := strings.Split(smsSvc.LastMessage, "\n")
		pincode := msg[0]
		hash2 := msg[len(msg)-1]

		assert.Equal(t, hash, hash2)

		res2, err := svc.VerifyOTP(ctx, otpsvc.VerifyOTPRequest{
			OTPUUID: res.OTPUUID,
			PinCode: pincode,
		})

		assert.NoError(t, err)
		assert.Equal(t, phoneNumber, res2.PhoneNumber)
	})
}
