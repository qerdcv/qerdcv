package repositories

import "errors"

var (
	ErrUniqueConstraint     = errors.New("unique constraint")
	ErrForeignKeyConstraint = errors.New("foreign key constraint")
	ErrNotFound             = errors.New("not found")
)
