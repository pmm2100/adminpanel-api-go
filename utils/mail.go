package utils

import (
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendMail(email string, msg []byte) error {
	godotenv.Load()

	smtp_conn := os.Getenv("SMTP_CONN")
	smtp_from := os.Getenv("SMTP_FROM")
	smtp_pass := os.Getenv("SMTP_PASS")
	smtp_host := os.Getenv("SMTP_HOST")
	smtp_user := os.Getenv("SMTP_USER")
	var emails []string
	emails = append(emails, email)

	auth := smtp.PlainAuth("", smtp_user, smtp_pass, smtp_host)

	err := smtp.SendMail(smtp_host+":"+smtp_conn, auth, smtp_from, emails, msg)
	if err != nil {
		return err
	}
	return nil
}
