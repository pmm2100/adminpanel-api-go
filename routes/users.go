package routes

import (
	"net/http"
	"strconv"

	"adminPanel.com/middlewares"
	"adminPanel.com/models"
	"adminPanel.com/utils"
	"github.com/gin-gonic/gin"
)

func createUser(context *gin.Context) {
	permBanned := middlewares.CheckForBannedIP(context)
	if permBanned {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You have been perm banned"})
		return
	}
	var user models.UserLogin
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Params!"})
		return
	}
	err = user.RegisterUser()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user, try again later!", "err": err.Error()})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Problem with generating JWT, try again later!"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Account created!", "token": token})
}

func loginUser(context *gin.Context) {
	permBanned := middlewares.CheckForBannedIP(context)
	if permBanned {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You have been perm banned"})
		return
	}
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Params!"})
		return
	}
	err = models.CheckLogIn(user.Username, user.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong username or Password!", "err": err.Error()})
		return
	}
	userData, err := models.GetUserBy("username", user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not read rest user data!"})
		return
	}

	token, err := utils.GenerateJWT(userData.ID, userData.IsAdmin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Problem with generating JWT, try again later!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User have logged in", "token": token})

}

func updateUser(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Params!"})
		return
	}
	userId := context.GetInt64("userId")
	err = user.UpdateUserData(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update user data"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Updating user data Succed"})
}

func deleteaccount(context *gin.Context) {
	userId := context.GetInt64("userId")
	err := models.DeleteUser(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Account deleted Succesfully"})

}

func getUserData(context *gin.Context) {

	var data models.Data
	err := context.ShouldBindJSON(&data)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Params!"})
		return
	}
	user, err := models.GetUserBy(data.SearchVal, data.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't get searched user data"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Data fetch succesfull", "data": user})
}

func getUsersData(context *gin.Context) {

	var data models.Data
	err := context.ShouldBindJSON(&data)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Params!"})
		return
	}
	users, err := models.GetUsersBy(data.SearchVal, data.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't get searched users data"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Data fetch succesfull", "data": users})
}

func banUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not convert context param to int"})
		return
	}
	user, err := models.GetUserBy("id", id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User with this id doesn't exist"})
		return
	}
	err = user.BanUser()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not ban user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Succesfully banned user"})

}

func adminDeleteUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not convert context param to int"})
		return
	}
	_, err = models.GetUserBy("id", id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User with this id doesn't exist"})
		return
	}
	err = models.AdminDeleteUser(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Succesfully deleted user"})
}

func adminUpdateUserData(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not convert context param to int"})
		return
	}
	_, err = models.GetUserBy("id", id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User with this id doesn't exist"})
		return
	}
	var user models.UserLogin
	err = context.ShouldBindJSON(&user)
	user.AdminUpdateUserData(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update user data"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Updating user data succesfull"})

}
