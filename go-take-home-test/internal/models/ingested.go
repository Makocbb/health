package models

import (
	"fmt"
	"time"
)

var ErrIngestedFormNotFound = fmt.Errorf("ingested form not found")

// IngestedForm is the schema currently agreed with the external provider.
type IngestedForm struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`

	SessionID            string          `json:"session_id" bun:"session_id"`
	ApplicationReference string          `json:"application_reference" bun:"application_reference"`
	Name                 string          `json:"name" bun:"name"`
	Email                string          `json:"email" bun:"email"`
	Gender               IngestedGender  `json:"gender" bun:"gender,type:jsonb"`
	DateOfBirth          string          `json:"date_of_birth" bun:"date_of_birth"`
	PhoneNumber          *string         `json:"phone_number" bun:"phone_number"`
	MobileNumber         string          `json:"mobile_number" bun:"mobile_number"`
	Address              IngestedAddress `json:"address" bun:"address,type:jsonb"`
}

type IngestedGender string

const (
	IngestedGenderMale   IngestedGender = "male"
	IngestedGenderFemale IngestedGender = "female"
	IngestedGenderOther  IngestedGender = "other"
)

type IngestedAddress struct {
	AddressLine1 string  `json:"address_line_1"`
	AddressLine2 string  `json:"address_line_2"`
	AddressLine3 *string `json:"address_line_3"`
	Postcode     string  `json:"postcode"`
	Country      string  `json:"country"`
}

type IngestedFormParams struct {
	Page                 int    `json:"page"`
	PerPage              int    `json:"per_page"`
	SessionID            string `json:"session_id"`
	ApplicationReference string `json:"application_reference"`
}
