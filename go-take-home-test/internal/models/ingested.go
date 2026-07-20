package models

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
