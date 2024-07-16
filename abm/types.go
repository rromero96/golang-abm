package abm

import "time"

const ABM_COLLECTION = "abm"

type DTO struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
