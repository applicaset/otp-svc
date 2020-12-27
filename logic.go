package otpsvc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/applicaset/sms-svc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math"
	"math/rand"
	"strconv"
	"text/template"
	"time"
)

type service struct {
	repo         Repository
	expiresInSec int
	rnd          *rand.Rand
	pinLen       int
	smsSvc       smssvc.Service
	msgTpl       *template.Template
}

func WithGoogleAutomaticSMSVerification(hash string) SendOTPOption {
	return func(data map[string]interface{}) {
		data["Hash"] = hash
	}
}

func (svc *service) SendOTP(ctx context.Context, req SendOTPRequest, options ...SendOTPOption) (*SendOTPResponse, error) {
	phoneNumber := req.PhoneNumber
	// TODO: validate phone number

	pinCode := svc.generatePinCode()

	entity := Entity{
		UUID:        uuid.New().String(),
		PhoneNumber: phoneNumber,
		PinCode:     pinCode,
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(svc.expiresInSec)),
	}

	err := svc.repo.Create(ctx, entity)
	if err != nil {
		return nil, errors.Wrap(err, "error on create entity")
	}

	data := map[string]interface{}{"PinCode": pinCode}

	for i := range options {
		options[i](data)
	}

	var buf bytes.Buffer

	err = svc.msgTpl.Execute(&buf, data)
	if err != nil {
		return nil, errors.Wrap(err, "error on execute message template")
	}

	err = svc.smsSvc.SendMessage(ctx, phoneNumber, buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "error on send message")
	}

	rsp := SendOTPResponse{OTPUUID: entity.UUID}

	return &rsp, nil
}

func (svc *service) generatePinCode() string {
	return fmt.Sprintf("%0"+strconv.Itoa(svc.pinLen)+"d", svc.rnd.Intn(int(math.Pow10(svc.pinLen))))
}

func (svc *service) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	entity, err := svc.repo.Find(ctx, req.OTPUUID)
	if err != nil {
		return nil, errors.Wrap(err, "error on find otp by uuid")
	}

	if entity == nil {
		return nil, ErrOTPNotFoundOrExpired{UUID: req.OTPUUID}
	}

	if entity.ExpiresAt.Before(time.Now()) {
		return nil, ErrOTPNotFoundOrExpired{UUID: req.OTPUUID}
	}

	if entity.PinCode != req.PinCode {
		return nil, ErrInvalidPinCode{UUID: req.OTPUUID, PinCode: req.PinCode}
	}

	rsp := VerifyOTPResponse{PhoneNumber: entity.PhoneNumber}

	return &rsp, nil
}

const defaultMessageTemplate = "{{ .PinCode }}"

const defaultPinLen = 4

const defaultExpiresInSec = 300

func New(repo Repository, smsSvc smssvc.Service, options ...Option) Service {
	svc := service{
		repo:         repo,
		expiresInSec: defaultExpiresInSec,
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
		pinLen:       defaultPinLen,
		smsSvc:       smsSvc,
		msgTpl:       template.Must(template.New("message").Parse(defaultMessageTemplate)),
	}

	for i := range options {
		options[i](&svc)
	}

	return &svc
}
