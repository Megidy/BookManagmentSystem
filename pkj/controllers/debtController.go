package controllers

import (
	"net/http"

	"github.com/Megidy/BookManagmentSystem/pkj/models"
	"github.com/Megidy/BookManagmentSystem/pkj/utils"
	"github.com/gin-gonic/gin"
)

// check all debts that users have	-admin
func CheckDebts(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.(*models.User).Role == "admin" {
		debts, err := models.GetAllDebts()
		utils.HandleError(c, err, "didnt retrieve any data from db", http.StatusInternalServerError)

		c.JSON(http.StatusOK, gin.H{
			"all debts": debts,
		})
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func CheckUsersDebt(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else {

		debts, err := models.GetAllUsersDebts(user.(*models.User))
		utils.HandleError(c, err, "didnt get any of debts from  db", http.StatusInternalServerError)

		c.JSON(http.StatusOK, gin.H{
			"your current debts:": debts,
		})

	}
}
