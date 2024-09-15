package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Megidy/BookManagmentSystem/pkj/models"
	"github.com/Megidy/BookManagmentSystem/pkj/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignUp(c *gin.Context) {
	var user models.User
	//{
	//	"username":"..."
	//	"password":"..."
	//}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}
	//better to add to database already created users with role admin
	//and not to be able to create new user with this role
	//but for simplicity i did like u can create user with role admin
	if user.Password == "AdminPass123" && user.Username == "SecretAdmin123" {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	NewUser := models.User{
		Username: user.Username,
		Password: string(hash),
		Role:     user.Role,
	}
	ok, err := models.IsSignedUp(&NewUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user is with this name is already signed up",
		})
		return
	}

	err = models.CreateUser(&NewUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

}
func LogIn(c *gin.Context) {
	var NewUserRequest models.User
	//{
	//	"username":"..."
	//	"password":"..."
	//}
	err := c.ShouldBindJSON(&NewUserRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	ok, err := models.IsSignedUp(&NewUserRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "didnt find user with this username",
		})
		return
	}
	user, err := models.FindUser(&NewUserRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	err = utils.UnHashPassword([]byte(user.Password), []byte(NewUserRequest.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*15, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"succescfull": "you just logged in",
	})
}
func DeleteUserAdmin(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if user.(*models.User).Role == "admin" {
		id := c.Param("userId")
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		user2, err := models.FindUserById(float64(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		if user2.Role == "admin" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cant delete admin",
			})
			return
		} else {
			err = models.DeleteUser(userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"successful": "account has been deleted",
			})

		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

}
func UpdateUserInfo(c *gin.Context) {
	var NewRequest struct {
		NewUsername string
		OldPassword string
		NewPassword string
	}
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err := c.ShouldBindJSON(&NewRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	if user.(*models.User).Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cant change admins info",
		})
		return
	} else {
		user2, err := models.FindUserById(float64(user.(*models.User).Id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		err = utils.UnHashPassword([]byte(user2.Password), []byte(NewRequest.OldPassword))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})
			return
		}
		hash, err := utils.HashPassword(NewRequest.NewPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to hash password",
			})
		}

		err = models.ChangeInfo(hash, user.(*models.User).Id, NewRequest.NewUsername)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"successful": "you changed your info",
		})

	}
}
func GetAllUser(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
	}
	if user.(*models.User).Role == "admin" {
		users, err := models.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
func DeleteAccount(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var NewUserRequest models.User
	err := c.ShouldBindJSON(&NewUserRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad Request",
		})
		return
	}
	usertemp, err := models.FindUserById(float64(user.(*models.User).Id))
	if usertemp.Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you can`t delete admin ((",
		})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	err = utils.UnHashPassword([]byte(usertemp.Password), []byte(NewUserRequest.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	err = models.DeleteUser(user.(*models.User).Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"successful": "you deleted your account",
	})

}
