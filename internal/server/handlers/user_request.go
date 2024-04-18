package handlers

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.Length(3, 35)),
		validation.Field(&r.Password, validation.Required, validation.Length(3, 35)),
	)
}

type AuthorizeUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r AuthorizeUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.Length(3, 35)),
		validation.Field(&r.Password, validation.Required, validation.Length(3, 35)),
	)
}
