package models

import (
	"errors"
	"fmt"

	"adminPanel.com/db"
	"adminPanel.com/utils"
)

type UserLogin struct {
	ID       int64
	Username string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
	IsAdmin  int
	IsBanned int
}

type User struct {
	ID       int64
	Username string `binding:"required"`
	Email    string
	Password string `binding:"required"`
	IsAdmin  int
	IsBanned int
}

type Data struct {
	SearchVal string `binding:"required"`
	Value     any    `binding:"required"`
}

func (u *UserLogin) RegisterUser() error {
	query := "INSERT INTO users (username, email, password) VALUES (?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Username, u.Email, hash)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	u.IsAdmin = 0
	u.IsBanned = 0
	return nil

}

func CheckLogIn(username string, password string) error {
	query := "SELECT id, password FROM users WHERE username=?"
	var hash string
	var id int
	row := db.DB.QueryRow(query, username)
	err := row.Scan(&id, &hash)
	if err != nil {
		return errors.New("Invalid Credentials")
	}
	err, banned := CheckForBanned(int64(id))
	if err != nil {
		return errors.New("Could not get banned data")
	}
	if banned == 1 {
		return errors.New("You have been banned")
	} else if banned != 1 && banned != 0 {
		return errors.New("Unexpected banned value")
	}
	ok := utils.ComparePasswordToHash(hash, password)
	if !ok {
		return errors.New("Invalid Password")
	}
	return nil

}

func GetUserBy(searchVal string, value any) (*User, error) {
	var u User

	allowed := map[string]bool{
		"id":       true,
		"username": true,
		"email":    true,
	}

	if !allowed[searchVal] {
		return nil, errors.New("invalid search field")
	}

	query := fmt.Sprintf(
		"SELECT id, username, email, password, isadmin, isbanned FROM users WHERE %v = ?",
		searchVal,
	)
	row := db.DB.QueryRow(query, value)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.IsAdmin, &u.IsBanned)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUsersBy(searchVal string, value any) ([]*User, error) {
	var userArr []*User

	query := fmt.Sprintf(
		"SELECT id, username, email, password, isadmin, isbanned FROM users WHERE %s LIKE ?",
		searchVal,
	)
	rows, err := db.DB.Query(query, fmt.Sprintf("%%%v%%", value))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.IsAdmin, &u.IsBanned)
		if err != nil {
			return nil, err
		}
		userArr = append(userArr, &u)
	}
	return userArr, nil
}

func (UpdatedUser *User) UpdateUserData(id int64) error {
	query := "UPDATE users SET username=?, password=? WHERE id = ?"
	hash, err := utils.HashPassword(UpdatedUser.Password)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(query, UpdatedUser.Username, hash, id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id int64) error {
	_, err := GetUserBy("id", id)
	if err != nil {
		return errors.New("User with this id doesn't exists")
	}
	query := "DELETE FROM users WHERE users.id=?"
	_, err = db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil

}

func AdminDeleteUser(id int64) error {
	_, err := GetUserBy("id", id)
	if err != nil {
		return errors.New("User with this id doesn't exists")
	}
	query := "DELETE FROM users WHERE id=?"
	_, err = db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (UpdatedUser *UserLogin) AdminUpdateUserData(id int64) error {
	query := "UPDATE users SET username=?, password=?, email=? WHERE id = ?"
	hash, err := utils.HashPassword(UpdatedUser.Password)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(query, UpdatedUser.Username, hash, UpdatedUser.Email, id)
	if err != nil {
		return err
	}
	return nil
}
