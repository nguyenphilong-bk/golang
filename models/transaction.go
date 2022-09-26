package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Type      string             `json:"type,omitempty"`
	Amount    int64              `json:"amount,omitempty"`
	Balance   int64              `json:"balance,omitempty"`
	From      string             `json:"from,omitempty"`
	To        string             `json:"to,omitempty"`
	CreatedAt int64              `json:"created_at,omitempty"`
	UpdatedAt int64              `json:"updated_at,omitempty"`
}

// TransactionModel ...
type TransactionModel struct{}

func (m TransactionModel) Create(ctx context.Context, form forms.CreateTransactionForm) (transaction Transaction, err error) {
	//Check if the user exists in database
	fmt.Println("Transaction model: Create")

	transactionCollection := db.GetCollection(db.DB, "transactions")

	result, err := transactionCollection.InsertOne(ctx, Transaction{
		Type:      form.Type,
		Amount:    form.Amount,
		Balance:   form.Balance,
		From:      form.From,
		To:        form.To,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})

	if err != nil {
		return transaction, errors.New("error when creating new transaction")
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID)
	transaction.Amount = form.Amount
	transaction.Type = form.Type
	transaction.Balance = form.Balance
	transaction.From = form.From
	transaction.To = form.To
	transaction.CreatedAt = form.CreatedAt
	transaction.UpdatedAt = form.UpdatedAt

	return transaction, err
}

func (m TransactionModel) Retrieve(ctx context.Context, user User) (transactions []Transaction, err error) {
	fmt.Println("Transaction model: Retrieve")

	transactionCollection := db.GetCollection(db.DB, "transactions")

	results, err := transactionCollection.Find(ctx, bson.M{"$or": []bson.M{
		{"from": user.Username},
		{"to": user.Username},
	}})

	if err != nil {
		return transactions, errors.New("error when retrieving transactions")
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var transaction Transaction
		if err = results.Decode(&transaction); err != nil {
			return transactions, errors.New("error when decoding transaction")
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
