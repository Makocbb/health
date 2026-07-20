package providers

import (
	"errors"
	"math/rand"
	"time"
)

var ErrSendFailed = errors.New("sendgrid: send failed")

type EmailRequest struct {
	To      string
	From    string
	Subject string
	Body    string
}

// SendEmail is a mock SendGrid client.
// It succeeds ~95% of the time and always sleeps ~1s to simulate latency.
func SendEmail(req EmailRequest) error {
	_ = req

	time.Sleep(time.Second)

	if rand.Float64() < 0.95 {
		return nil
	}

	return ErrSendFailed
}
