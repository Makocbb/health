package models

import (
	"fmt"
	"time"
)

var ErrTransformLogNotFound = fmt.Errorf("transform log not found")

type TransformLog struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`

	Success bool   `json:"success" bun:"success"`
	Error   string `json:"message" bun:"message"`

	TransformedFormID int64 `json:"transformed_form_id" bun:"transformed_form_id"`
}

type TransformLogParams struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Success *bool `json:"success"`
}
