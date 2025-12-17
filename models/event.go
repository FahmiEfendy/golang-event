package models

import (
	"fmt"
	"time"

	"example.com/event/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

func (e Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, datetime, user_id) 
	VALUES (?, ?, ?, ?, ?)
	`

	// DB.Prepare is used to create a prepared statement for execution
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// Ensure the statement is closed after execution
	defer stmt.Close()

	// stmt.Exec is used to execute a prepared statement with the given arguments
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	// Get the last inserted ID and assign it to the event
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println("Last inserted ID is", id)

	return nil
}

func GetAllEvents() ([]Event, error) {
	query := `
	SELECT * FROM events
	`

	// DB.Query is used to execute a query that returns rows
	rows, err := db.DB.Query(query)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return nil, err
	}

	// defer rows.Close() ensures that the rows are closed after processing
	defer rows.Close()

	// Slice to hold the retrieved events
	events := []Event{}
	for rows.Next() {
		var e Event

		err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func GetEventByID(eventId int64) (*Event, error) {
	query := `
	SELECT * FROM events WHERE id = ?
	`

	// QueryRow is used to execute a query that is expected to return at most one row
	row := db.DB.QueryRow(query, eventId)

	var e Event

	err := row.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events 
	SET name = ?, description = ?, location = ?, datetime = ? 
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// Ensure the statement is closed after execution
	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (event Event) Delete() error {
	query := `
	DELETE FROM events WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// Ensure the statement is closed after execution
	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	if err != nil {
		return err
	}

	return nil
}
