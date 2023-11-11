package authdb

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

type User struct {
	ID       int    `json: "user_id"`
	Name     string `json:"user_name"`
	Password string `json:"user_password"`
}

func Connect(dbRoot string, dbPassword string, dbHost string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/", dbRoot, dbPassword, dbHost))
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

func CreateDB(db *sql.DB) {
	cmd, err := db.Query("CREATE DATABASE IF NOT EXISTS auth")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer cmd.Close()
}
func CreateTables(db *sql.DB) {
	cmd, err := db.Query("CREATE TABLE IF NOT EXISTS auth.users (user_id int AUTO_INCREMENT,  user_name char(50) NOT NULL, user_password char(128), PRIMARY KEY(user_id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer cmd.Close()
}
func InsertUser(db *sql.DB, user User) error {
	password := md5.Sum([]byte(user.Password))
	cmd, err := db.Query(fmt.Sprintf("INSERT INTO auth.users (user_name,user_password) VALUES ('%s','%s');", user.Name, hex.EncodeToString(password[:])))
	if err != nil {
		return err
	}
	defer cmd.Close()
	return nil
}
func GetUserByName(user_name string, db *sql.DB) (User, error) {
	var user User
	results, err := db.Query(fmt.Sprintf("SELECT * FROM auth.users where user_name = '%s'", user_name))
	if err != nil {
		return user, err
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Password)
		if err != nil {
			return user, err
		}
	}
	return user, nil
}
func CreateUser(db *sql.DB, u User) (bool, error) {
	user, err := GetUserByName(u.Name, db)
	if err != nil {
		return false, err
	}
	if user != (User{}) {
		return false, nil
	} else {
		err := InsertUser(db, u)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}
