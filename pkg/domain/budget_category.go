package domain

import "github.com/google/uuid"

type Category struct {
	ID     uuid.UUID
	UserID int
	Name   string
}
