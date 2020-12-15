package otp_svc

import "context"

type Service interface {
	SendOTP(ctx context.Context, req SendOTPRequest) (res *SendOTPResponse, err error)
	VerifyOTP(ctx context.Context, req VerifyOTPRequest) (res *VerifyOTPResponse, err error)
}

type SendOTPRequest struct {
	PhoneNumber string
}

type SendOTPResponse struct {
	OTPUUID string
}

type VerifyOTPRequest struct {
	OTPUUID string
	PinCode string
}

type VerifyOTPResponse struct {
	PhoneNumber string
}
