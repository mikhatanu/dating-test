package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/mikhatanu/dating-test/db"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Id       string
	Username string
	Password string
}

// Get user using sql
func GetUser(username string) (user, error) {
	// Create sql query template
	get := `select rowid, username, password from user where username = '%s';`

	// execute query
	result := db.DB.QueryRow(fmt.Sprintf(get, username))

	// save the result to struct
	user := user{}
	err := result.Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	}

	return user, nil

}

// sign up user using sql
func Signup(username string, password string) error {
	if username == "" || password == "" {
		return errors.New("username or password is missing")
	}
	// hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error when generating hashed password")
	}

	// Create sql query template
	insert := `insert into user (username, password) VALUES ('%s','%s');`

	// execute query
	_, err2 := db.DB.Exec(fmt.Sprintf(insert, username, hashedPassword))
	if err2 != nil {
		log.Println(err2.Error())
		return errors.New("error occured when creating user")
	}
	log.Println("Success creating user")

	return nil
}

// check password
func CheckPassword(password, hashedPassword string) error {
	// check password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return errors.New("incorrect password")
	}
	return nil
}
