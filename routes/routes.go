package routes

import (
	"adminPanel.com/middlewares"
	"github.com/gin-gonic/gin"
)

func ChechForRoutes(server *gin.Engine) {

	server.POST("/register", createUser)
	server.POST("/login", loginUser)

	server.POST("/register/admin", createAdmin)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.DELETE("/delete", deleteaccount)
	authenticated.PUT("/update", updateUser)

	authenticatedAdmin := server.Group("/")
	authenticatedAdmin.Use(middlewares.AuthenticateAdmin)

	authenticatedAdmin.GET("/searchUser", getUserData)
	authenticatedAdmin.GET("/searchUsers", getUsersData)
	authenticatedAdmin.PUT("/ban/:id", banUser)
	authenticatedAdmin.DELETE("/deleteuser/:id", adminDeleteUser)
	authenticatedAdmin.PUT("/updateadmin/:id", adminUpdateUserData)
	authenticatedAdmin.PUT("/permban/:id", permBanUser)
	authenticatedAdmin.GET("/getPermbanned", getPermBannedIPs)
	authenticatedAdmin.PUT("/revokeban/:id", revokeBan)
	authenticatedAdmin.DELETE("/revokePermban", revokePermBan)

	server.GET("/sendVeryfication/:mail", sendMail)
}
