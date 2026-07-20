package models

import (
	"fmt"
	"strings"
	"time"
)

var (
	ErrIngestedFormNotFound  = fmt.Errorf("ingested form not found")
	ErrIngestedFormDuplicate = fmt.Errorf("ingested form already exists")
	ErrInvalidDateOfBirth    = fmt.Errorf("invalid date of birth")
	ErrInvalidIngestedForm   = fmt.Errorf("invalid ingested form")
)

const (
	IngestedStatusPending     = "pending"
	IngestedStatusTransformed = "transformed"
	IngestedStatusFailed      = "failed"
)

// IngestedForm is the schema currently agreed with the external provider.
type IngestedForm struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`

	Status string `json:"status" bun:"status"`

	Fingerprint string `json:"fingerprint" bun:"fingerprint"`

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
		IngestedFormID:       i.ID,
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

// Validate checks required fields against the currently agreed ingest schema.
// Extra JSON fields are ignored by encoding/json (tolerant of provider drift).
func (i *IngestedForm) Validate() error {
	switch {
	case strings.TrimSpace(i.SessionID) == "":
		return fmt.Errorf("%w: session_id is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.ApplicationReference) == "":
		return fmt.Errorf("%w: application_reference is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.Name) == "":
		return fmt.Errorf("%w: name is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.Email) == "":
		return fmt.Errorf("%w: email is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(string(i.Gender)) == "":
		return fmt.Errorf("%w: gender is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.DateOfBirth) == "":
		return fmt.Errorf("%w: date_of_birth is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.MobileNumber) == "":
		return fmt.Errorf("%w: mobile_number is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.Address.AddressLine1) == "":
		return fmt.Errorf("%w: address.address_line_1 is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.Address.AddressLine2) == "":
		return fmt.Errorf("%w: address.address_line_2 is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.Address.Postcode) == "":
		return fmt.Errorf("%w: address.postcode is required", ErrInvalidIngestedForm)
	case strings.TrimSpace(i.Address.Country) == "":
		return fmt.Errorf("%w: address.country is required", ErrInvalidIngestedForm)
	}

	switch i.Gender {
	case IngestedGenderMale, IngestedGenderFemale, IngestedGenderOther:
	default:
		return fmt.Errorf("%w: gender must be male, female, or other", ErrInvalidIngestedForm)
	}

	if _, err := time.Parse("2006-01-02", i.DateOfBirth); err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidDateOfBirth, i.DateOfBirth)
	}
	return nil
}
