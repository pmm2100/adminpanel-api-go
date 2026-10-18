package routes

import (
	"net/http"

	"adminPanel.com/models"
	"adminPanel.com/utils"
	"github.com/gin-gonic/gin"
)

func createAdmin(context *gin.Context) {
	var user models.UserLogin
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Params!"})
		return
	}
	err = user.RegisterAdmin()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user, try again later!"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Problem with generating JWT, try again later!"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Account created!", "token": token})
}
