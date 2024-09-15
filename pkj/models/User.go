package models

import (
	"database/sql"
	"fmt"

	"github.com/Megidy/BookManagmentSystem/pkj/config"
)

var db *sql.DB

func init() {
	config.Connect()
	db = config.GetDb()
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func CreateUser(user *User) error {
	_, err := db.Exec("insert into users(username,password,role) values(?,?,?)",
		user.Username, user.Password, user.Role)
	if err != nil {
		fmt.Println("didnt create new user")
		return err
	}
	return nil

}

func IsSignedUp(NewUser *User) (bool, error) {
	var user User
	row := db.QueryRow("select username from users where username =? ",
		NewUser.Username)
	err := row.Scan(&user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, meaning the user is not signed up
			return false, nil
		} else {
			return true, err
		}

	}
	if user.Username == NewUser.Username {
		return true, nil
	}
	return false, nil
}

func FindUser(user *User) (*User, error) {
	var SignedUser User
	row := db.QueryRow("select * from users where username =?", user.Username)
	err := row.Scan(&SignedUser.Id, &SignedUser.Username, &SignedUser.Password, &SignedUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return &User{}, err
		}

	}
	return &SignedUser, err
}
func FindUserById(id float64) (*User, error) {
	var user User
	row := db.QueryRow("select * from users where id = ?", id)

	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
	}

	return &user, nil
}

func DeleteUser(id int) error {

	_, err := db.Exec("delete from users where id=?", id)
	if err != nil {
		return err
	}
	return nil
}
func GetAllUsers() (*[]User, error) {
	var Users []User
	query, err := db.Query("select * from users")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	for query.Next() {
		var user User
		err := query.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
		Users = append(Users, user)

	}
	return &Users, nil
}
func ChangeInfo(password []byte, id int, username string) error {
	_, err := db.Exec("update users set password= ?,username=? where id =?", password, username, id)
	if err != nil {
		return err
	}
	return nil
}
