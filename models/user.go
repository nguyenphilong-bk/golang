package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Username  string             `json:"username,omitempty"`
	Name      string             `json:"name,omitempty"`
	Password  string             `json:"-"`
	UpdatedAt int64              `json:"updated_at,omitempty"`
	CreatedAt int64              `json:"created_at,omitempty"`
	Balance   int64              `json:"balance,omitempty"`
}

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)
var transactionModel = new(TransactionModel)

// Login ...
func (m UserModel) Login(form forms.LoginForm) (user User, token Token, err error) {
	fmt.Println("User model: Login")
	userCollection := db.GetCollection(db.DB, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"username": form.Username}).Decode(&user)
	if err != nil {
		return user, token, err
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, token, err
	}

	//Generate the JWT auth token
	tokenDetails, err := authModel.CreateToken(user.ID.Hex())
	if err != nil {
		return user, token, err
	}

	saveErr := authModel.CreateAuth(user.ID.Hex(), tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	return user, token, nil
}

// Register ...
func (m UserModel) Register(form forms.RegisterForm) (user User, err error) {
	//Check if the user exists in database
	fmt.Println("User model: Register")

	userCollection := db.GetCollection(db.DB, "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.D{{
		Key:   "username",
		Value: form.Username,
	}}).Decode(&user)

	if err != nil && err != mongo.ErrNoDocuments {
		return user, errors.New("something went wrong, please try again later")
	}

	if err == mongo.ErrNoDocuments {
		bytePassword := []byte(form.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
		if err != nil {
			return user, errors.New("something went wrong, please try again later")
		}
		newUser := User{
			ID:        primitive.NewObjectID(),
			Username:  form.Username,
			Password:  string(hashedPassword),
			Name:      form.Name,
			UpdatedAt: time.Now().Unix(),
			CreatedAt: time.Now().Unix(),
			Balance:   0,
		}
		_, insertError := userCollection.InsertOne(ctx, newUser)

		if insertError != nil {
			return user, errors.New("error when inserting new user")
		}

		user.ID = newUser.ID
		user.Name = newUser.Name
		user.Username = newUser.Username

		return user, err
	}

	return user, errors.New("username already existed")
}

func (m UserModel) TopUp(ctx context.Context, userID primitive.ObjectID, form forms.TopUpForm) (transaction Transaction, err error) {
	//Check if the user exists in database
	fmt.Println("User model: TopUp")

	userCollection := db.GetCollection(db.GetDB(), "users")

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := db.DB.StartSession()
	if err != nil {
		return transaction, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var user User
		err = userCollection.FindOne(sessionContext, bson.M{"id": userID}).Decode(&user)
		if err != nil && err != mongo.ErrNoDocuments {
			return transaction, errors.New("something went wrong, please try again later")
		}

		now := time.Now().Unix()

		update := bson.M{"balance": user.Balance + form.Amount, "updatedat": now}
		result, err := userCollection.UpdateOne(sessionContext, bson.M{"id": userID}, bson.M{"$set": update})

		if err != nil {
			return transaction, errors.New("internal server error")
		}

		var updatedUser User

		if result.MatchedCount == 1 {
			err = userCollection.FindOne(sessionContext, bson.M{"id": userID}).Decode(&updatedUser)
		}

		if err != nil {
			return transaction, errors.New("internal server error")
		}

		transaction, err = transactionModel.Create(ctx, forms.CreateTransactionForm{
			From:      updatedUser.Username,
			To:        updatedUser.Username,
			Amount:    form.Amount,
			Balance:   updatedUser.Balance,
			Type:      utils.TOP_UP,
			CreatedAt: now,
			UpdatedAt: now,
		})

		return transaction, err
	}

	data, err := session.WithTransaction(ctx, callback, txnOpts)
	v, _ := data.(Transaction)

	return v, err
}

func (m UserModel) WithDraw(ctx context.Context, userID primitive.ObjectID, form forms.WithDrawForm) (transaction Transaction, err error) {
	//Check if the user exists in database
	fmt.Println("User model: WithDraw")

	userCollection := db.GetCollection(db.DB, "users")
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := db.DB.StartSession()
	if err != nil {
		return transaction, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var user User
		err = userCollection.FindOne(sessionContext, bson.M{"id": userID}).Decode(&user)
		if err != nil && err != mongo.ErrNoDocuments {
			return user, errors.New("something went wrong, please try again later")
		}

		if user.Balance < form.Amount {
			return user, errors.New("your balance is not enough to withdraw")
		}

		now := time.Now().Unix()

		update := bson.M{"balance": user.Balance - form.Amount, "updatedat": now}
		result, err := userCollection.UpdateOne(sessionContext, bson.M{"id": userID}, bson.M{"$set": update})

		if err != nil {
			return user, errors.New("internal server error")
		}

		//get updated user details
		var updatedUser User
		if result.MatchedCount == 1 {
			err = userCollection.FindOne(sessionContext, bson.M{"id": userID}).Decode(&updatedUser)
		}

		if err != nil {
			return user, errors.New("internal server error")
		}

		transaction, err = transactionModel.Create(sessionContext, forms.CreateTransactionForm{
			From:      updatedUser.Username,
			To:        updatedUser.Username,
			Amount:    form.Amount,
			Balance:   updatedUser.Balance,
			Type:      utils.WITHDRAW,
			CreatedAt: now,
			UpdatedAt: now,
		})

		return transaction, err
	}
	data, err := session.WithTransaction(ctx, callback, txnOpts)
	v, _ := data.(Transaction)

	return v, err
}

func (m UserModel) Details(userId primitive.ObjectID, ctx context.Context) (transactions []Transaction, err error) {
	fmt.Println("User model: Details")
	userCollection := db.GetCollection(db.DB, "users")
	var user User

	err = userCollection.FindOne(ctx, bson.M{"id": userId}).Decode(&user)
	if err != nil {
		return transactions, err
	}

	transactions, err = transactionModel.Retrieve(ctx, user)

	return transactions, err
}

func (m UserModel) Transfer(ctx context.Context, userId primitive.ObjectID, form forms.TransferForm) (transaction Transaction, err error) {
	fmt.Println("User model: Transfer")
	userCollection := db.GetCollection(db.DB, "users")

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := db.DB.StartSession()
	if err != nil {
		return transaction, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var source, target User
		err = userCollection.FindOne(ctx, bson.M{"id": userId}).Decode(&source)
		if err != nil {
			return transaction, err
		}

		err = userCollection.FindOne(ctx, bson.M{"username": form.To}).Decode(&target)

		if err != nil {
			return transaction, errors.New("target user not existed")
		}

		if source.Balance < form.Amount {
			return transaction, errors.New("your balance is not enough to execute the transaction")
		}
		now := time.Now().Unix()

		sourceUpdate, err := userCollection.UpdateOne(sessionContext, bson.M{"id": source.ID}, bson.M{"$set": bson.M{"balance": source.Balance - form.Amount, "updatedat": now}})

		if err != nil {
			return sourceUpdate, errors.New("internal server error")
		}

		targetUpdate, err := userCollection.UpdateOne(sessionContext, bson.M{"id": target.ID}, bson.M{"$set": bson.M{"balance": target.Balance + form.Amount, "updatedat": now}})

		if err != nil {
			return targetUpdate, errors.New("internal server error")
		}
		
		transaction, err = transactionModel.Create(sessionContext, forms.CreateTransactionForm{
			From:      source.Username,
			To:        form.To,
			Amount:    form.Amount,
			Balance:   source.Balance,
			Type:      utils.TRANSFER,
			CreatedAt: now,
			UpdatedAt: now,
		})
		
		return transaction, err
	}
	data, err := session.WithTransaction(ctx, callback, txnOpts)
	value, _ := data.(Transaction)

	return value, err
}