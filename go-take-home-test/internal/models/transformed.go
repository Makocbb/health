package models

import (
	"fmt"
	"time"
)

var ErrTransformedFormNotFound = fmt.Errorf("transformed form not found")

// TransformedForm is the schema expected by the FORM-BOT after processing.
type TransformedForm struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`

	SentToBot bool `json:"sent_to_bot" bun:"sent_to_bot"`

	IngestedFormID       int64             `json:"ingested_form_id" bun:"ingested_form_id"`
	SessionID            string            `json:"session_id" bun:"session_id"`
	ApplicationReference string            `json:"application_reference" bun:"application_reference"`
	FirstName            string            `json:"first_name" bun:"first_name"`
	LastName             string            `json:"last_name" bun:"last_name"`
	Email                string            `json:"email" bun:"email"`
	Gender               TransformedGender `json:"gender" bun:"gender"`
	DateOfBirth          time.Time         `json:"date_of_birth" bun:"date_of_birth"`
	PhoneNumber          *string           `json:"phone_number" bun:"phone_number"`
	MobileNumber         string            `json:"mobile_number" bun:"mobile_number"`
	AddressLine1         string            `json:"address_line_1" bun:"address_line_1"`
	AddressLine2         string            `json:"address_line_2" bun:"address_line_2"`
	AddressLine3         *string           `json:"address_line_3" bun:"address_line_3"`
	Postcode             string            `json:"postcode" bun:"postcode"`
	Country              string            `json:"country" bun:"country"`
	Longitude            float64           `json:"longitude" bun:"longitude"`
	Latitude             float64           `json:"latitude" bun:"latitude"`
}

type TransformedGender string

const (
	TransformedGenderMale           TransformedGender = "male"
	TransformedGenderFemale         TransformedGender = "female"
	TransformedGenderPreferNotToSay TransformedGender = "prefer-not-to-say"
)

type TransformedFormParams struct {
	Page                 int    `json:"page"`
	PerPage              int    `json:"per_page"`
	SessionID            string `json:"session_id"`
	ApplicationReference string `json:"application_reference"`
}
