package req

import "github.com/go-playground/validator/v10"

type UserReq struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateReq struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,min=3,max=32"`
	Password string `json:"password"`
}

func (u *UserReq) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserUpdateReq) ValidateOnUpdate() error {
	validate := validator.New()
	return validate.Struct(u)
}
