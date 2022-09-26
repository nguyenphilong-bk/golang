package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryParamInt(c *gin.Context, name string, defaultValue int) (int, error) {
  param := c.Query(name)
  result, err := strconv.Atoi(param)
  if err != nil {
    return defaultValue, errors.New(name + " param can not be parsed to integer")
  }
  return result, nil
}