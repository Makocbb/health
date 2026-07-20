package presenters

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"go-take-home-test/internal/models"
)

// IngestedForm is the schema currently agreed with the external provider.
type IngestedForm struct {
	SessionID            string          `json:"session_id"`
	ApplicationReference string          `json:"application_reference"`
	Name                 string          `json:"name"`
	Email                string          `json:"email"`
	Gender               IngestedGender  `json:"gender"`
	DateOfBirth          string          `json:"date_of_birth"`
	PhoneNumber          *string         `json:"phone_number"`
	MobileNumber         string          `json:"mobile_number"`
	Address              IngestedAddress `json:"address"`
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

// Fingerprint returns a stable hash of the entire ingest payload contents.
func (i *IngestedForm) Fingerprint() (string, error) {
	payload, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func (i *IngestedForm) ToModel() (*models.IngestedForm, error) {
	fp, err := i.Fingerprint()
	if err != nil {
		return nil, err
	}
	return &models.IngestedForm{
		Fingerprint:          fp,
		SessionID:            i.SessionID,
		ApplicationReference: i.ApplicationReference,
		Name:                 i.Name,
		Email:                i.Email,
		Gender:               models.IngestedGender(i.Gender),
		DateOfBirth:          i.DateOfBirth,
		PhoneNumber:          i.PhoneNumber,
		MobileNumber:         i.MobileNumber,
		Address: models.IngestedAddress{
			AddressLine1: i.Address.AddressLine1,
			AddressLine2: i.Address.AddressLine2,
			AddressLine3: i.Address.AddressLine3,
			Postcode:     i.Address.Postcode,
			Country:      i.Address.Country,
		},
	}, nil
}
