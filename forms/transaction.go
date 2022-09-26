package forms

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type TransactionForm struct{}

type CreateTransactionForm struct {
	From      string `form:"from" json:"from,omitempty"`
	To        string `form:"to" json:"to,omitempty"`
	Amount    int64  `from:"amount" json:"amount,omitempty" binding:"required,min=0"`
	Balance   int64  `from:"balance" json:"balance,omitempty" binding:"required,min=0"`
	Type      string `form:"type" json:"type,omitempty" binding:"required"`
	CreatedAt int64  `form:"created_at" json:"created_at,omitempty"`
	UpdatedAt int64  `form:"updated_at" json:"updated_at,omitempty"`
}

type TransferForm struct {
	To string `form:"to" json:"to,omitempty"`
	Amount int64 `form:"amount" json:"amount,omitempty" binding:"required,min=0"`
}

func (f TransactionForm) From(tag string, errMsg ...string) string {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the source account"
		}
		return errMsg[0]
	case "min", "max":
		return "The source account must have 5 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f TransactionForm) To(tag string, errMsg ...string) string {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the target account"
		}
		return errMsg[0]
	case "min", "max":
		return "The target account must have 5 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f TransactionForm) Amount(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Amount can't be blank or equal to 0"
		}
		return errMsg[0]
	case "min":
		return "The amount must greater than 0"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f TransactionForm) Balance(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Balance can't be blank"
		}
		return errMsg[0]
	case "min":
		return "The balance must greater than 0"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f TransactionForm) Type(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Type can't be blank"
		}
		return errMsg[0]
	default:
		return "Something went wrong, please try again later"
	}
}

func (f TransactionForm) Transfer(err error) string {
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

			// if err.Field() == "Balance" {
			// 	return f.Balance(err.Tag())
			// }

			// if err.Field() == "From" {
			// 	return f.From(err.Tag())
			// }

			if err.Field() == "To" {
				return f.To(err.Tag())
			}

			// if err.Field() == "Type" {
			// 	return f.Type(err.Tag())
			// }
		}

	default:
		return "Invalid payload"
	}

	return "Something went wrong, please try again later"
}

