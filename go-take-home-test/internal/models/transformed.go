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

	SessionID            string            `json:"sessionId" bun:"session_id"`
	ApplicationReference string            `json:"applicationReference" bun:"application_reference"`
	FirstName            string            `json:"firstName" bun:"first_name"`
	LastName             string            `json:"lastName" bun:"last_name"`
	Email                string            `json:"email" bun:"email"`
	Gender               TransformedGender `json:"gender" bun:"gender,type:jsonb"`
	DateOfBirth          time.Time         `json:"dateOfBirth" bun:"date_of_birth,type:jsonb"`
	PhoneNumber          *string           `json:"phoneNumber" bun:"phone_number"`
	MobileNumber         string            `json:"mobileNumber" bun:"mobile_number"`
	AddressLine1         string            `json:"addressLine1" bun:"address_line_1"`
	AddressLine2         string            `json:"addressLine2" bun:"address_line_2"`
	AddressLine3         *string           `json:"addressLine3" bun:"address_line_3"`
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
