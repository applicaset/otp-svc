package otp_svc

import "time"

type Entity struct {
	UUID        string
	PhoneNumber string
	PinCode     string
	ExpiresAt   time.Time
}
