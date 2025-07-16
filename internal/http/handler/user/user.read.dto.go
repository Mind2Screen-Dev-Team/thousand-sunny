package user

import (
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
	"github.com/google/uuid"
)

type (
	ExampleUserReadRequestInput struct {
		ID uuid.UUID `path:"id" example:"a4c55779-44d3-469b-ac15-a28401e9b80f" format:"uuid" doc:"Unique identifier of the user paramater" required:"true"`
	}

	ExampleUserReadResponseOutput struct {
		Body   ExampleUserReadResponseBody
		Status int
	}
)

type (
	ExampleUserReadResponseData struct {
		ID        uuid.UUID `json:"id" doc:"Unique identifier of the user" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
		Name      string    `json:"name" doc:"Full name of the user" example:"John Doe"`
		Age       int       `json:"age" doc:"Age of the user in years" example:"30"`
		CreatedAt time.Time `json:"created_at" doc:"Timestamp when the user was created" example:"2024-07-16T15:04:05Z" format:"date-time"`
		UpdatedAt time.Time `json:"updated_at" doc:"Timestamp when the user was last updated" example:"2024-07-16T16:30:00Z" format:"date-time"`
	}
	ExampleUserReadResponseBody xresp.GeneralResponse[*ExampleUserReadResponseData, any]
)
