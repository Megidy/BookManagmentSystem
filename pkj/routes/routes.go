package routes

import (
	"github.com/Megidy/BookManagmentSystem/pkj/controllers"
	"github.com/Megidy/BookManagmentSystem/pkj/middlware"
	"github.com/gin-gonic/gin"
)

var InitRoutes = func(router gin.IRouter) {
	//user auth
	router.POST("/signup/", controllers.SignUp)
	router.POST("/login/", controllers.LogIn)

	//for user role=admin
	router.POST("/admin/books/", middlware.RequireAuth, controllers.CreateBook)
	router.DELETE("admin/books/:bookId/", middlware.RequireAuth, controllers.DeleteBook)
	router.PUT("/admin/books/:bookId/", middlware.RequireAuth, controllers.UpdateBook)
	router.GET("/admin/debts/", middlware.RequireAuth, controllers.CheckDebts)
	router.DELETE("/admin/users/:userId/", middlware.RequireAuth, controllers.DeleteUserAdmin)
	router.GET("/admin/users/", middlware.RequireAuth, controllers.GetAllUser)
	//for user role=users
	router.GET("/books/", middlware.RequireAuth, controllers.GetAllBooks)
	router.GET("/books/:bookId/", middlware.RequireAuth, controllers.GetBookById)
	router.PUT("/givebackbook/", middlware.RequireAuth, controllers.TakeBook)
	router.PUT("/takebook/", middlware.RequireAuth, controllers.GiveBook)
	router.GET("/mydebts/", middlware.RequireAuth, controllers.CheckUsersDebt)
	router.PUT("/myaccount/alter/", middlware.RequireAuth, controllers.UpdateUserInfo)
	router.DELETE("/myaccount/delete/", middlware.RequireAuth, controllers.DeleteAccount)
}
