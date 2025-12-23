package models

import "example.com/event/db"

func (event Event) RegisterEvent(userId int64) error {
	query := `
		INSERT INTO registrations(event_id, user_id) VALUES (?, ?)
	`

	// DB.Prepare is used to create a prepared statement for execution
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// Ensure the statement is closed after execution
	defer stmt.Close()

	// stmt.Exec is used to execute a prepared statement with the given arguments
	_, err = stmt.Exec(event.ID, userId)
	if err != nil {
		return err
	}

	return nil
}

func (event Event) UnregisterEvent(userId int64) error {
	query := `
	DELETE FROM registrations WHERE event_id = ? AND user_id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// Ensure the statement is closed after execution
	defer stmt.Close()

	_, err = stmt.Exec(event.ID, userId)
	if err != nil {
		return err
	}

	return nil
}
