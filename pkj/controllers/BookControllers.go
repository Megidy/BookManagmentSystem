package controllers

import (
	"net/http"
	"strconv"

	"github.com/Megidy/BookManagmentSystem/pkj/models"
	"github.com/Megidy/BookManagmentSystem/pkj/utils"
	"github.com/gin-gonic/gin"
)

// adding new type of book to database	- admin
func CreateBook(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
	if user.(*models.User).Role == "admin" {
		var NewBookRequest models.Book

		c.ShouldBindJSON(&NewBookRequest)
		ok, err := models.IsCreated(NewBookRequest)

		utils.HandleError(c, err, "something went wrong when checking book in db", http.StatusInternalServerError)

		if ok {

			utils.HandleError(c, nil, "book is already created", http.StatusBadRequest)

			err = models.CreateBook(&NewBookRequest)
			utils.HandleError(c, err, "", http.StatusInternalServerError)
			c.JSON(http.StatusOK, gin.H{
				"succesfully created new book": NewBookRequest,
			})

		} else {
			c.AbortWithStatus(http.StatusNotFound)

		}
	}
}

// check all avaible books	-admin,user
func GetAllBooks(c *gin.Context) {
	books, err := models.GetAllBooks()
	utils.HandleError(c, err, "didnt retrieve books from db", http.StatusInternalServerError)
	c.JSON(http.StatusOK, gin.H{
		"all books ": books,
	})
}

// check book with current id	-admin,user
func GetBookById(c *gin.Context) {
	id := c.Param("bookId")
	bookId, err := strconv.Atoi(id)

	utils.HandleError(c, err, "didnt converte to int ", http.StatusBadRequest)

	book, err := models.FindBook(bookId)
	utils.HandleError(c, err, "didnt find book in db", http.StatusInternalServerError)

	c.JSON(http.StatusOK, gin.H{
		"book": book,
	})
}

// delete data about book from database	-admin
func DeleteBook(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
	if user.(*models.User).Role == "admin" {
		id := c.Param("bookId")
		bookId, err := strconv.Atoi(id)
		utils.HandleError(c, err, "didnt converte to int ", http.StatusBadRequest)

		book, err := models.FindBook(bookId)
		utils.HandleError(c, err, "didnt find book in db", http.StatusInternalServerError)

		err = models.DeleteBook(bookId)
		utils.HandleError(c, err, "didnt delete book from db", http.StatusInternalServerError)

		c.JSON(http.StatusOK, gin.H{
			"deleted book": book,
		})

	} else {
		c.AbortWithStatus(http.StatusNotFound)

	}

}

// Updating data about book	-admin
func UpdateBook(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
	if user.(*models.User).Role == "admin" {
		var NewBookRequest models.Book

		c.ShouldBindJSON(&NewBookRequest)
		oldbook, err := models.FindBook(NewBookRequest.Id)
		utils.HandleError(c, err, "didnt find book in db", http.StatusInternalServerError)
		err = models.UpdateBook(&NewBookRequest)
		utils.HandleError(c, err, "didnt update book in db", http.StatusInternalServerError)

		c.JSON(http.StatusOK, gin.H{

			"after":  NewBookRequest,
			"before": oldbook,
		})
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

// +1 to avaible_copies	-admin,user
func TakeBook(c *gin.Context) {
	var debt models.Debt
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var NewRequest models.Debt

	c.ShouldBindJSON(&NewRequest)
	if NewRequest.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "amount<=0",
		})
		return
	}
	debt.Book_id = NewRequest.Book_id
	debt.Amount = NewRequest.Amount
	ok, err := models.AlreadyExist(user.(*models.User), &debt)
	utils.HandleError(c, err, "something went wron during cheching book avaibility in db", http.StatusInternalServerError)

	//if true than compare input data and debt datafrom db, reset if needed
	if ok {
		debt2, err := models.GetDebt(user.(*models.User), &debt)
		utils.HandleError(c, err, "didnt get debt from db", http.StatusInternalServerError)

		//if amount request is < to amount of debt than error
		if debt2.Amount-debt.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "u didnt take that much books",
			})
			return
		}
		//if amount request is == to amount of debt than reseting debt
		if debt2.Amount-debt.Amount == 0 && debt2.Book_id == debt.Book_id {
			err = models.TakeBook(debt.Book_id, debt.Amount)
			utils.HandleError(c, err, "didnt retrieve book from  db", http.StatusInternalServerError)

			err := models.ResetDebt(user.(*models.User), &debt)
			utils.HandleError(c, err, "didnt reset debt in db", http.StatusInternalServerError)
			//if amount request is > to amount of debt than updating debt
		} else if debt2.Amount-debt.Amount > 0 && debt2.Book_id == debt.Book_id {
			if debt2.Amount-debt.Amount > 10 {
				c.JSON(http.StatusBadRequest, gin.H{
					"sorry": "u cant take more than 10 books once",
				})
			}
			err := models.TakeBook(debt.Book_id, debt2.Amount)
			utils.HandleError(c, err, "didnt update amouint of books in db", http.StatusInternalServerError)

			err = models.SubAmountFromDebt(user.(*models.User), &debt)
			utils.HandleError(c, err, "didnt sub from book amount in db", http.StatusInternalServerError)

			c.JSON(http.StatusOK, gin.H{})
		}

		//no data if didnt find debt
	} else {
		c.JSON(http.StatusOK, gin.H{
			"no data about debt": nil,
		})
		return
	}

}

// -1 to avaible_copies	-admin,user
func GiveBook(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
	var NewRequest models.Debt

	c.ShouldBindJSON(&NewRequest)
	if NewRequest.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "amount<=0",
		})
		return
	}
	//firstly decrementing amount of books from database and than creating or updating new debt
	//in query checking if amount of request  is > than amount in databaase
	err := models.GiveBook(NewRequest.Book_id, NewRequest.Amount)
	utils.HandleError(c, err, "didnt take book from db", http.StatusInternalServerError)

	var debt models.Debt
	debt.Book_id = NewRequest.Book_id
	debt.Amount = NewRequest.Amount
	//checking if debt is already exist
	ok, err = models.AlreadyExist(user.(*models.User), &debt)
	utils.HandleError(c, err, "didnt check existence of book in db", http.StatusInternalServerError)

	//if debt doesnt exist create new
	if !ok {

		err = models.CreateDebt(user.(*models.User).Id, &debt)
		utils.HandleError(c, err, "didnt create new debt in db", http.StatusInternalServerError)

		c.JSON(http.StatusOK, nil)

	} else {
		// if debt already exist than update this
		err := models.AddAmountForDebt(user.(*models.User), &debt)
		utils.HandleError(c, err, "didnt add amount of books in db", http.StatusInternalServerError)

		//
		c.JSON(http.StatusOK, nil)
	}

}
