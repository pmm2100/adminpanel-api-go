package models

import (
	"database/sql"
	"errors"

	"adminPanel.com/db"
)

type PermBan struct {
	ID int
	IP string
}

func (u *User) PermBan(IP string) error {
	u.BanUser()
	query := "INSERT INTO permBanned (userIP) VALUES (?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(IP)
	if err != nil {
		return err
	}
	return nil
}

func GetPermBannedIPs() (error, []PermBan) {
	var Permbanned []PermBan
	query := "SELECT * FROM permBanned"
	rows, err := db.DB.Query(query)
	if err == sql.ErrNoRows {
		Permbanned = make([]PermBan, 0)
		return nil, Permbanned
	} else if err != nil {
		return err, nil
	}

	for rows.Next() {
		var Permban PermBan
		err := rows.Scan(&Permban.ID, &Permban.IP)
		if err != nil {
			return err, nil
		}
		Permbanned = append(Permbanned, Permban)
	}
	return nil, Permbanned
}

func RevokePermBanIP(IP string) error {
	query := "DELETE FROM permBanned WHERE userIP=?"
	_, err := db.DB.Exec(query, IP)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CheckForAdmin() (int, error) {
	query := "SELECT isadmin FROM users WHERE id=?"
	row := db.DB.QueryRow(query, u.ID)
	err := row.Scan(&u.IsAdmin)
	if err != nil {
		return 0, err
	}
	return u.IsAdmin, nil
}

func (u *User) BanUser() error {
	_, err := GetUserBy("id", u.ID)
	if err != nil {
		return errors.New("User with this id doesn't exists")
	}
	query := "UPDATE users SET isbanned=1 WHERE id=?"
	_, err = db.DB.Exec(query, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func RevokeBan(id int64) error {
	_, err := GetUserBy("id", id)
	if err != nil {
		return errors.New("User with this id doesn't exists")
	}
	query := "UPDATE users SET isbanned=0 WHERE id=?"
	_, err = db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func CheckForBanned(id int64) (error, int) {
	query := "SELECT isbanned FROM users WHERE id=?"
	var banned int
	row := db.DB.QueryRow(query, id)
	row.Scan(&banned)
	switch banned {
	case 1:
		return nil, 1
	case 2:
		return nil, 0
	default:
		return errors.New("Unexpected value of isbanned"), 0
	}
}
