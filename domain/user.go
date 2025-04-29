package domain

import (
	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID    uuid.UUID `bson:"id" json:"id"`
	Name  string    `bson:"name" json:"name"`
	Email string    `bson:"email" json:"email"`
}
