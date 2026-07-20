package models

import "errors"

var ErrSendFailed = errors.New("sendgrid: send failed")

type EmailRequest struct {
	To      string
	From    string
	Subject string
	Body    string
}
