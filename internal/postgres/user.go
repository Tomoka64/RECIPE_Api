package postgres

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *User) CreateUser(exec Executer) error {
	statement := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := exec.Exec(statement, u.Name, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(exec Executer, name, password string) (*User, error) {
	var user *User
	if err := exec.Get(user, "SELECT * FROM users WHERE name = $1 and password = $2", name, password); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func UserExistsByUsername(exec Executer, name string) bool {
	var user *User
	return exec.Get(user, "SELECT * FROM users WHERE name = $1", name) == nil
}

func UserExistsByEmail(exec Executer, email string) bool {
	var user *User
	return exec.Get(user, "SELECT * FROM users WHERE email = $1", email) == nil
}
