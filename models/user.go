package models

import (
	"errors"

	"example.com/event/db"
	"example.com/event/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required,email"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users (email, password) VALUES (?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	query := `
		SELECT id, password FROM users WHERE email = ?
	`

	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string

	// Scan will return error if no row matches the query
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		// Return error if no user found
		return err
	}

	isPasswordValid := utils.CheckPassword(u.Password, retrievedPassword)

	if !isPasswordValid {
		return errors.New("invalid credentials")
	}

	return nil
}
