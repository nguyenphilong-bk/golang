package forms

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// UserForm ...
type UserForm struct{}

// LoginForm ...
type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required,min=5,max=5"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

// RegisterForm ...
type RegisterForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=20,fullName"` //fullName rule is in validator.go
	Username string `form:"username" json:"username" binding:"required,min=5,max=5"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

type TopUpForm struct {
	Amount int64 `form:"amount" json:"amount" binding:"min=0,required"`
}

type WithDrawForm struct {
	Amount int64 `form:"amount" json:"amount" binding:"min=0,required"`
}

// Name ...
func (f UserForm) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your name"
		}
		return errMsg[0]
	case "min", "max":
		return "Your name should be between 3 to 20 characters"
	case "fullName":
		return "Name should not include any special characters or numbers"
	default:
		return "Something went wrong, please try again later"
	}
}

// Username ...
func (f UserForm) Username(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your email"
		}
		return errMsg[0]
	case "min", "max":
		return "Username is just 5 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

// Password ...
func (f UserForm) Password(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your password"
	case "min", "max":
		return "Your password should be between 3 and 50 characters"
	case "eqfield":
		return "Your passwords does not match"
	default:
		return "Something went wrong, please try again later"
	}
}

// Signin ...
func (f UserForm) Login(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Username" {
				return f.Username(err.Tag())
			}
			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

// Register ...
func (f UserForm) Register(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return f.Name(err.Tag())
			}

			if err.Field() == "Username" {
				return f.Username(err.Tag())
			}

			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}

		}
	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

// Amount ...
func (f UserForm) Amount(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the amount of money that you want to top-up"
		}
		return errMsg[0]
	case "min":
		return "The amount must greater than 0"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f UserForm) TopUp(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("%v", ok)
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Amount" {
				return f.Amount(err.Tag())
			}
		}

	default:
		return "Invalid payload"
	}

	return "Something went wrong, please try again later"
}

func (f UserForm) WithDraw(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("%v", ok)
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Amount" {
				return f.Amount(err.Tag())
			}
		}

	default:
		return "Invalid payload"
	}

	return "Something went wrong, please try again later"
}

