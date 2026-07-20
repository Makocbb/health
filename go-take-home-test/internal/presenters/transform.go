package presenters

import (
	"go-take-home-test/internal/models"
	"time"
)

type TransformRequest struct {
	IngestedID int64 `json:"ingested_id"`
}

type TransformedFormBotRequest struct {
	SessionID            string                   `json:"sessionId" bun:"session_id"`
	ApplicationReference string                   `json:"applicationReference" bun:"application_reference"`
	FirstName            string                   `json:"firstName" bun:"first_name"`
	LastName             string                   `json:"lastName" bun:"last_name"`
	Email                string                   `json:"email" bun:"email"`
	Gender               models.TransformedGender `json:"gender" bun:"gender,type:jsonb"`
	DateOfBirth          time.Time                `json:"dateOfBirth" bun:"date_of_birth,type:jsonb"`
	PhoneNumber          *string                  `json:"phoneNumber" bun:"phone_number"`
	MobileNumber         string                   `json:"mobileNumber" bun:"mobile_number"`
	AddressLine1         string                   `json:"addressLine1" bun:"address_line_1"`
	AddressLine2         string                   `json:"addressLine2" bun:"address_line_2"`
	AddressLine3         *string                  `json:"addressLine3" bun:"address_line_3"`
	Postcode             string                   `json:"postcode" bun:"postcode"`
	Country              string                   `json:"country" bun:"country"`
	Longitude            float64                  `json:"longitude" bun:"longitude"`
	Latitude             float64                  `json:"latitude" bun:"latitude"`
}

func TransformedFromModel(model *models.TransformedForm) *TransformedFormBotRequest {
	return &TransformedFormBotRequest{
		SessionID:            model.SessionID,
		ApplicationReference: model.ApplicationReference,
		FirstName:            model.FirstName,
		LastName:             model.LastName,
		Email:                model.Email,
		Gender:               model.Gender,
		DateOfBirth:          model.DateOfBirth,
		PhoneNumber:          model.PhoneNumber,
		MobileNumber:         model.MobileNumber,
		AddressLine1:         model.AddressLine1,
		AddressLine2:         model.AddressLine2,
		AddressLine3:         model.AddressLine3,
		Postcode:             model.Postcode,
		Country:              model.Country,
		Longitude:            model.Longitude,
		Latitude:             model.Latitude,
	}
}

type SendToBotRequest struct {
	TransformedID int64 `json:"transformed_id"`
}
