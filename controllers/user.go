package controllers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/Massad/gin-boilerplate/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController ...
type UserController struct{}

var userModel = new(models.UserModel)
var userForm = new(forms.UserForm)
var transactionForm = new(forms.TransactionForm)

// getUserID ...
func getUserID(c *gin.Context) (userID primitive.ObjectID) {
	return c.MustGet("userID").(primitive.ObjectID)
}

// Login ...

// @Summary Login api
// @Schemes
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/login [post]
// @Param username body string true "username" SchemaExample(Subject: longn)
// @Param password body string true "password" SchemaExample(Subject: malongnhan)
func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		message := userForm.Login(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: message})
		return
	}

	user, token, err := userModel.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{Status: http.StatusUnauthorized, Message: "The username or password is incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

// Register ...
// @Summary Register api
// @Schemes
// @Description Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/register [post]
// @Param name body string true "name" SchemaExample(Subject: long nguyen)
// @Param username body string true "username" SchemaExample(Subject: longn)
// @Param password body string true "password" SchemaExample(Subject: malongnhan)
func (ctrl UserController) Register(c *gin.Context) {
	var registerForm forms.RegisterForm
	if validationErr := c.ShouldBindJSON(&registerForm); validationErr != nil {
		message := userForm.Register(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: message})
		return
	}

	if !utils.IsStringAlphabetic(registerForm.Username) {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: "The username must not contain any special characters"})
		return 
	} 

	user, err := userModel.Register(registerForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusNotAcceptable, Message: err.Error()})
		return
	}

	temp, _ := json.Marshal(&user)
	var result map[string]interface{}
	json.Unmarshal(temp, &result)

	c.JSON(http.StatusOK, utils.Response{Status: http.StatusOK, Message: "Register new account successfully", Data: result})
}

// @Summary Top-up api
// @Schemes
// @Description Top-up to my account
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/top-up [post]
// @Param amount body int true "amount of money" SchemaExample(S/tranubject: 5000)
func (ctrl UserController) TopUp(c *gin.Context) {
	userID := getUserID(c)

	var form forms.TopUpForm
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := userForm.TopUp(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: message})
		return
	}

	transaction, err := userModel.TopUp(ctx, userID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	temp, _ := json.Marshal(&transaction)
	var result map[string]interface{}
	json.Unmarshal(temp, &result)
	c.JSON(http.StatusOK, utils.Response{Status: http.StatusOK, Message: "Top-up successfully", Data: result})
}

// @Summary Withdraw api
// @Schemes
// @Description Withdraw from my account
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/withdraw [post]
// @Param amount body int true "username of target account" SchemaExample(Subject: 5000)
func (ctrl UserController) WithDraw(c *gin.Context) {
	userID := getUserID(c)

	var form forms.WithDrawForm
	ctx, cancel := context.WithTimeout(context.Background(), 10*time .Second)
	defer cancel()

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := userForm.WithDraw(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: message})
		return
	}

	transaction, err := userModel.WithDraw(ctx, userID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	temp, _ := json.Marshal(&transaction)
	var result map[string]interface{}
	json.Unmarshal(temp, &result)
	c.JSON(http.StatusOK, utils.Response{Status: http.StatusOK, Message: "Withdraw successfully", Data: result})
}

// @Summary Details api
// @Schemes
// @Description Get details transactions from my account
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/details [get]
func (ctrl UserController) Details(c *gin.Context) {
	userID := getUserID(c)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
		
	// transactions, err := userModel.Details(userID, ctx, query)
	transactions, err := userModel.Details(userID, ctx)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Response{Status: http.StatusNotAcceptable, Message: err.Error(), Data: nil})
		return
	}

	data := make([]interface{}, len(transactions))
	for i, v := range transactions {
		data[i] = v
	}

	c.JSON(http.StatusOK, utils.RetrieveResponse{Status: http.StatusOK, Message: "Retrieve user details successfully", Data: data})
}

// @Summary Transfer api
// @Schemes
// @Description Transfer to another account
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/transfer [post]
// @Param to body string true "Target account" SchemaExample(longn)
// @Param amount body int true "Amount of money" SchemaExample(5000)
func (ctrl UserController) Transfer(c *gin.Context) {
	userID := getUserID(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var form forms.TransferForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := transactionForm.Transfer(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Response{Status: http.StatusNotAcceptable, Message: message})
		return
	}

	transaction, err := userModel.Transfer(ctx, userID, form)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Response{Status: http.StatusNotAcceptable, Message: err.Error(), Data: nil})
		return
	}

	temp, _ := json.Marshal(&transaction)
	var result map[string]interface{}
	json.Unmarshal(temp, &result)

	c.JSON(http.StatusOK, utils.Response{Status: http.StatusOK, Message: "Transaction created successfully", Data: result})
}
