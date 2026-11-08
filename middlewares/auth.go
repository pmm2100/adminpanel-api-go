package middlewares

import (
	"database/sql"
	"net/http"

	"adminPanel.com/db"
	"adminPanel.com/models"
	"adminPanel.com/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	permBanned := CheckForBannedIP(context)
	if permBanned {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Your IP have been perm banned"})
		return
	}
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}
	id, _, err := utils.CheckJWT(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}
	err, banned := models.CheckForBanned(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get banned value"})
		return
	}
	if banned == 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You have been banned!"})
		return
	}
	context.Set("userId", id)
	context.Next()

}

func AuthenticateAdmin(context *gin.Context) {
	permBanned := CheckForBannedIP(context)
	if permBanned {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Your IP have been perm banned"})
		return
	}
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}
	id, isAdmin, err := utils.CheckJWT(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}
	if isAdmin != 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You don't have enough priviledges for this action"})
		return
	}
	err, banned := models.CheckForBanned(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get banned value"})
		return
	}
	if banned == 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You have been banned!"})
		return
	}
	context.Set("userId", id)
	context.Next()
}

func CheckForBannedIP(context *gin.Context) bool {
	IP := context.ClientIP()
	var bannedIP string
	query := "SELECT userIP FROM permBanned WHERE userIP=?"
	err := db.DB.QueryRow(query, IP).Scan(&bannedIP)
	if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}

}

func AuthenticatePasswordReset(context *gin.Context) {
	permBanned := CheckForBannedIP(context)
	if permBanned {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Your IP have been perm banned"})
		return
	}
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}
	id, _, err := utils.CheckJWT(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}
	err, banned := models.CheckForBanned(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get banned value"})
		return
	}
	if banned == 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You have been banned!"})
		return
	}
	user, err := models.GetUserBy("id", id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get user data"})
		return
	}
	verified, err := utils.VerifyPasswordReset(user.Email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not check for "})
		return
	}
	switch verified {
	case false:
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Password reset session is not up to date"})
		return
	case true:
		context.Set("userId", id)
		context.Next()
	default:
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

}
