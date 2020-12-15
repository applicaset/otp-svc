package otp_svc

import "fmt"

type ErrOTPNotFoundOrExpired struct {
	UUID string
}

func (err ErrOTPNotFoundOrExpired) Error() string {
	return fmt.Sprintf("otp with uuid '%s' not found or expired", err.UUID)
}

type ErrInvalidPinCode struct {
	UUID    string
	PinCode string
}

func (err ErrInvalidPinCode) Error() string {
	return fmt.Sprintf("pin code '%s' is not valid for otp with uuid '%s'", err.PinCode, err.UUID)
}
