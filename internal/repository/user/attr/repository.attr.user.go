package attr

import (
	"time"

	"github.com/google/uuid"
)

type (
	ExampleUser struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Age       int       `json:"age"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
