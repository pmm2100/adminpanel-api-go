package models

import (
	"adminPanel.com/db"
	"adminPanel.com/utils"
)

func (u *UserLogin) RegisterAdmin() error {
	query := "INSERT INTO users (username, email, password, isadmin) VALUES (?,?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Username, u.Email, hash, 1)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	u.IsAdmin = 1
	u.IsBanned = 0
	return nil

}
