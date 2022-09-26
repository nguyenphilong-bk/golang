package controllers

import (
	"net/http"

	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-gonic/gin"
)

//AuthController ...
type AuthController struct{}

var authModel = new(models.AuthModel)

//TokenValid ...
func (ctl AuthController) TokenValid(c *gin.Context) {

	tokenAuth, err := authModel.ExtractTokenMetadata(c.Request)

	if err != nil {
		//Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	userID, err := authModel.FetchAuth(tokenAuth)
	if err != nil {
		//Token does not exists in Redis (User logged out or expired)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	//To be called from GetUserID()
	c.Set("userID", userID)
}
