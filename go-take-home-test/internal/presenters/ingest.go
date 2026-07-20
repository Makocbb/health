package presenters

import "go-take-home-test/internal/models"

// IngestedForm is the schema currently agreed with the external provider.
type IngestedForm struct {
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

func (i *IngestedForm) ToModel() *models.IngestedForm {
	return &models.IngestedForm{
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
	}
}
