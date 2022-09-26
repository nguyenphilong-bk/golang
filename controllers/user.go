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

// @Summary login api
// @Schemes
// @Description login
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Success"
// @Router /v1/user/login [post]
// @Param forms.LoginForm
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

func (ctrl UserController) WithDraw(c *gin.Context) {
	userID := getUserID(c)

	var form forms.WithDrawForm
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (ctrl UserController) Details(c *gin.Context) {
	userID := getUserID(c)
	
	// page, err := utils.QueryParamInt(c, "page", 1)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Response{Status: http.StatusNotAcceptable, Message: err.Error(), Data: nil}) 
	// 	return
	// }
	
	// // limit, err := utils.QueryParamInt(c, "limit", 20)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Response{Status: http.StatusNotAcceptable, Message: err.Error(), Data: nil}) 
	// 	return
	// }

	// query := models.Query{
	// 	Page: page,
	// 	Limit: limit,
	// 	Order: c.DefaultQuery("order", "desc"),
	// }

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
