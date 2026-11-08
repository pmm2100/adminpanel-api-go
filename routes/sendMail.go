package routes

import (
	"net/http"

	"adminPanel.com/models"
	"adminPanel.com/utils"

	"github.com/gin-gonic/gin"
)

func sendPasswordResetEmail(context *gin.Context) {
	email := context.Param("email")
	err := utils.SendPasswordResetEmail(email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not send password reset email", "err": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Password reset email send check your mailbox"})
}

func verifyUser(context *gin.Context) {
	type verifyEmail struct {
		Email      string `binding:"required"`
		Code       string `binding:"required"`
		ActionType string `binding:"required"`
	}
	var verifyData verifyEmail
	context.ShouldBindJSON(&verifyData)

	success, err := utils.VerifyCode(verifyData.Email, verifyData.Code, verifyData.ActionType)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check code", "err": err.Error()})
		return
	}
	if success {
		user, err := models.GetUserBy("email", verifyData.Email)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get data for this user"})
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Problem with generating JWT, try again later!"})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": "Account verified succesfully!", "token": token})

	} else if !success {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong code or the code is Expired"})
	}

}
