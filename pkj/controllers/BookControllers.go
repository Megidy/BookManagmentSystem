package controllers

import (
	"net/http"
	"strconv"

	"github.com/Megidy/BookManagmentSystem/pkj/models"
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
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "book is already created",
			})
			return
		}
		err = models.CreateBook(&NewBookRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"succesfully created new book": NewBookRequest,
		})

	} else {
		c.AbortWithStatus(http.StatusNotFound)

	}
}

// check all avaible books	-admin,user
func GetAllBooks(c *gin.Context) {
	books, err := models.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"all books ": books,
	})
}

// check book with current id	-admin,user
func GetBookById(c *gin.Context) {
	id := c.Param("bookId")
	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}
	book, err := models.FindBook(bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
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
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
		}
		book, err := models.FindBook(bookId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
		err = models.DeleteBook(bookId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
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
		err := models.UpdateBook(&NewBookRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
		book := NewBookRequest
		c.JSON(http.StatusOK, gin.H{
			"before": book,
			"after":  NewBookRequest,
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
	//{
	//"book_id":"..."
	//"amount:"..."
	//}
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	//if true than compare input data and debt datafrom db, reset if needed
	if ok {
		debt2, err := models.GetDebt(user.(*models.User), &debt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
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
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error ": err,
				})
				return
			}
			err := models.ResetDebt(user.(*models.User), &debt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
			//if amount request is > to amount of debt than updating debt
		} else if debt2.Amount-debt.Amount > 0 && debt2.Book_id == debt.Book_id {
			err := models.TakeBook(debt.Book_id, debt2.Amount)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
			err = models.SubAmountFromDebt(user.(*models.User), &debt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Error": err,
				})
				return
			}
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
		return
	}
	var NewRequest models.Debt
	//{
	//"book_id":"..."
	//"amount:"..."
	//}
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	var debt models.Debt
	debt.Book_id = NewRequest.Book_id
	debt.Amount = NewRequest.Amount
	//checking if debt is already exist
	ok, err = models.AlreadyExist(user.(*models.User), &debt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	//if debt doesnt exist create new
	if !ok {

		err = models.CreateDebt(user.(*models.User).Id, &debt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, nil)

	} else {
		// if debt already exist than update this
		err := models.AddAmountForDebt(user.(*models.User), &debt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		//
		c.JSON(http.StatusOK, nil)
	}

}
