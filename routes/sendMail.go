package routes

import (
	"fmt"
	"net/http"

	"adminPanel.com/utils"

	"github.com/gin-gonic/gin"
)

func sendMail(context *gin.Context) {
	email := string(context.Param("mail"))

	headerSubject := "Subject: Testowy Temat\r\n"
	headerTo := fmt.Sprintf("To: %s\r\n", email)

	body := "To jest treść wiadomości."

	msg := []byte(headerSubject + headerTo + body)

	err := utils.SendMail(email, msg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not send mail", "err": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Mail send, check spam"})
}
