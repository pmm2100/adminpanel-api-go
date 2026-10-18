package routes

import (
	"net/http"
	"strconv"

	"adminPanel.com/models"
	"github.com/gin-gonic/gin"
)

func permBanUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not convert context param to int"})
		return
	}
	u, err := models.GetUserBy("id", id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
		return
	}
	IP := context.ClientIP()
	err = u.PermBan(IP)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update db"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "IP and the user has been permbanned succesfully"})
}

func getPermBannedIPs(context *gin.Context) {
	Permbanned, err := models.GetPermBannedIPs()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get permbanned IP's"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Succesfully got permbanned IP's", "permbanned": Permbanned})
}

func revokeBan(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not convert context param to int"})
		return
	}
	err = models.RevokeBan(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not revoke ban"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Ban revoked Succefully"})
}

func revokePermBan(context *gin.Context) {
	var PermBan models.PermBan
	err := context.ShouldBindJSON(&PermBan)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not convert context param to int"})
		return
	}
	err = models.RevokePermBanIP(PermBan.IP)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not revoke permban"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Permban revoked Succefully"})
}
