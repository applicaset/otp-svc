package otpsvc

import "text/template"

type Option func(*service)

func WithMessageTemplate(tpl *template.Template) Option {
	return func(svc *service) {
		svc.msgTpl = tpl
	}
}

func WithGoogleAutomaticSMSVerificationTemplate() Option {
	tpl := `{{ .PinCode }}{{ if ne ( index . "Hash") "" }}


{{ .Hash }}{{ end }}`

	return func(svc *service) {
		svc.msgTpl = template.Must(template.New("message").Parse(tpl))
	}
}

func WithExpiresIn(sec int) Option {
	return func(svc *service) {
		svc.expiresInSec = sec
	}
}

func WithPinCodeLength(l int) Option {
	return func(svc *service) {
		svc.pinLen = l
	}
}
