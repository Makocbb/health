package models

import "errors"

var ErrSendFailed = errors.New("sendgrid: send failed")

type EmailRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    []byte `json:"body"`
}
