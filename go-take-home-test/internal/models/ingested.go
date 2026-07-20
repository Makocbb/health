package models

import (
	"fmt"
	"strings"
	"time"
)

var (
	ErrIngestedFormNotFound = fmt.Errorf("ingested form not found")
	ErrInvalidDateOfBirth   = fmt.Errorf("invalid date of birth")
)

// IngestedForm is the schema currently agreed with the external provider.
type IngestedForm struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`

	SessionID            string          `json:"session_id" bun:"session_id"`
	ApplicationReference string          `json:"application_reference" bun:"application_reference"`
	Name                 string          `json:"name" bun:"name"`
	Email                string          `json:"email" bun:"email"`
	Gender               IngestedGender  `json:"gender" bun:"gender"`
	DateOfBirth          string          `json:"date_of_birth" bun:"date_of_birth"`
	PhoneNumber          *string         `json:"phone_number" bun:"phone_number"`
	MobileNumber         string          `json:"mobile_number" bun:"mobile_number"`
	Address              IngestedAddress `json:"address" bun:"address,type:json"`
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

// ToTransformedForm maps an ingested form into the FORM-BOT schema.
// Name is split on the first space; gender "other" becomes "prefer-not-to-say".
// Coordinates should come from a prior postcode lookup.
func (i *IngestedForm) ToTransformedForm() (*TransformedForm, error) {
	dob, err := time.Parse("2006-01-02", i.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrInvalidDateOfBirth, i.DateOfBirth)
	}

	firstName, lastName := splitName(i.Name)

	form := &TransformedForm{
		SessionID:            i.SessionID,
		ApplicationReference: i.ApplicationReference,
		FirstName:            firstName,
		LastName:             lastName,
		Email:                i.Email,
		Gender:               i.Gender.ToTransformed(),
		DateOfBirth:          dob,
		PhoneNumber:          i.PhoneNumber,
		MobileNumber:         i.MobileNumber,
		AddressLine1:         i.Address.AddressLine1,
		AddressLine2:         i.Address.AddressLine2,
		AddressLine3:         i.Address.AddressLine3,
		Postcode:             i.Address.Postcode,
		Country:              i.Address.Country,
	}

	return form, nil
}

func (g IngestedGender) ToTransformed() TransformedGender {
	switch g {
	case IngestedGenderMale:
		return TransformedGenderMale
	case IngestedGenderFemale:
		return TransformedGenderFemale
	case IngestedGenderOther:
		return TransformedGenderPreferNotToSay
	default:
		return TransformedGenderPreferNotToSay
	}
}

func splitName(name string) (firstName, lastName string) {
	parts := strings.Fields(strings.TrimSpace(name))
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], strings.Join(parts[1:], " ")
	}
}
