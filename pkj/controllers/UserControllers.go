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
	utils.HandleError(c, err, "failed to hash password", http.StatusBadRequest)

	NewUser := models.User{
		Username: user.Username,
		Password: string(hash),
		Role:     user.Role,
	}
	ok, err := models.IsSignedUp(&NewUser)
	utils.HandleError(c, err, "didnt check if user signed up ", http.StatusBadRequest)

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user is with this name is already signed up",
		})
		c.Abort()
	}

	err = models.CreateUser(&NewUser)
	utils.HandleError(c, err, "failed to crate user", http.StatusBadRequest)

}
func LogIn(c *gin.Context) {
	var NewUserRequest models.User
	//{
	//	"username":"..."
	//	"password":"..."
	//}
	err := c.ShouldBindJSON(&NewUserRequest)
	utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

	ok, err := models.IsSignedUp(&NewUserRequest)
	utils.HandleError(c, err, "didnt find signed user ", http.StatusBadRequest)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "didnt find user with this username",
		})
		c.Abort()
	}
	user, err := models.FindUser(&NewUserRequest)
	utils.HandleError(c, err, "didnt find user in db", http.StatusBadRequest)

	err = utils.UnHashPassword([]byte(user.Password), []byte(NewUserRequest.Password))
	utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

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
		utils.HandleError(c, err, "failed to converte to int ", http.StatusBadRequest)

		user2, err := models.FindUserById(float64(userId))
		utils.HandleError(c, err, "didnt find user in db", http.StatusInternalServerError)

		if user2.Role == "admin" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cant delete admin",
			})
			c.Abort()
		} else {
			err = models.DeleteUser(userId)
			utils.HandleError(c, err, "didnt delete user in db", http.StatusInternalServerError)

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
	utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

	if user.(*models.User).Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cant change admins info",
		})
		c.Abort()
	} else {
		user2, err := models.FindUserById(float64(user.(*models.User).Id))
		utils.HandleError(c, err, "didnt find user ", http.StatusInternalServerError)

		err = utils.UnHashPassword([]byte(user2.Password), []byte(NewRequest.OldPassword))
		utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

		hash, err := utils.HashPassword(NewRequest.NewPassword)
		utils.HandleError(c, err, "failed to hash password", http.StatusBadRequest)

		err = models.ChangeInfo(hash, user.(*models.User).Id, NewRequest.NewUsername)
		utils.HandleError(c, err, "didnt change info ", http.StatusInternalServerError)

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
		utils.HandleError(c, err, "didnt retrieve users from db ", http.StatusInternalServerError)

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
	utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

	usertemp, err := models.FindUserById(float64(user.(*models.User).Id))
	if usertemp.Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you can`t delete admin ((",
		})
		c.Abort()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		c.Abort()
	}
	err = utils.UnHashPassword([]byte(usertemp.Password), []byte(NewUserRequest.Password))
	utils.HandleError(c, err, "Bad request", http.StatusBadRequest)

	err = models.DeleteUser(user.(*models.User).Id)
	utils.HandleError(c, err, "didnt delete user from db", http.StatusInternalServerError)

	c.JSON(http.StatusOK, gin.H{
		"successful": "you deleted your account",
	})

}
