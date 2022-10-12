// Package Routes provides the necessary functionality for roting
package routes

import (
	"app/db"
	"app/users"

	"github.com/julienschmidt/httprouter"
)

var Routes = func(router *httprouter.Router) {
	user := users.NewUserController(db.New())
	router.GET("/public", user.Public)
	router.GET("/user", user.Profile)
	router.POST("/signup", user.Signup)
	router.POST("/login", user.Login)
	router.GET("/search/:user", user.Search)
	router.POST("/upload", user.UploadImage)
	router.DELETE("/image/:id", user.DeleteImage)
	router.DELETE("/user", user.DeleteUser)
	router.PATCH("/user", user.UpdateInfo)
	router.GET("/logout", user.Logout)
}
