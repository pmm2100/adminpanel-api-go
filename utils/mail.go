package utils

import (
	"database/sql"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"adminPanel.com/db"
	"github.com/joho/godotenv"
)

func SendAuthenitcationEmail(email string) error {

	headerSubject := "Subject: Account veryfication \r\n"
	headerTo := fmt.Sprintf("To: %s\r\n\r\n", email)

	code := GenerateEmailAuthenticationCode()

	body := fmt.Sprintf("Your verify code: %v .", code)

	msg := []byte(headerSubject + headerTo + body)

	err := SendMail(string(email), []byte(msg))
	if err != nil {
		return err
	}

	date := time.Now().UTC().Add(15 * time.Minute)
	date_format := date.Format("2006-01-02 15:04:05")

	query := "INSERT OR REPLACE INTO verify (email, code, action_type, date) VALUES (?,?,?,?)"

	_, err = db.DB.Exec(query, email, code, "REGISTER", date_format)
	if err != nil {
		return err
	}
	return nil

}

func SendPasswordResetEmail(email string) error {
	headerSubject := "Subject: Password Reset\r\n"
	headerTo := fmt.Sprintf("To: %s\r\n\r\n", email)
	link := "" // add link to website here
	body := fmt.Sprintf("To reset your password click this link %v", link)

	msg := []byte(headerSubject + headerTo + body)

	err := SendMail(string(email), []byte(msg))
	if err != nil {
		return err
	}
	date := time.Now().UTC().Add(15 * time.Minute)
	date_format := date.Format("2006-01-02 15:04:05")
	query := "INSERT INTO verify (email, action_type, date) VALUES (?,?,?)"
	_, err = db.DB.Exec(query, email, "PASSWORDRESET", date_format)
	if err != nil {
		return err
	}

	return nil
}

func VerifyPasswordReset(email string) (bool, error) {
	var dbEmail string
	query := "SELECT email FROM verify WHERE email=? AND action_type=?"
	err := db.DB.QueryRow(query, email, "PASSWORDRESET").Scan(&dbEmail)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func VerifyCode(email, code, actionType string) (bool, error) {
	var dbEmail string
	query := `
		SELECT email FROM verify
		WHERE email = ? AND code = ? AND action_type = ? AND date >= datetime('now');
	`
	err := db.DB.QueryRow(query, email, code, actionType).Scan(&dbEmail)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	_, err = db.DB.Exec("DELETE FROM verify WHERE email = ? AND action_type = ?", email, actionType)
	if err != nil {
		return false, err
	}
	_, err = db.DB.Exec("UPDATE users SET verified=1 WHERE users.email=?", email)
	if err != nil {
		return false, err
	}
	return true, nil
}

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
