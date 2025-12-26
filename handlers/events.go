package handlers

import "example.com/event/models"

var GetAllEvents = func() ([]models.Event, error) {
	return models.GetAllEvents()
}

var GetEventByID = func(eventId int64) (*models.Event, error) {
	return models.GetEventByID(eventId)
}

var CreateEvent = func(event *models.Event) error {
	return event.Save()
}

var UpdateEvent = func(event *models.Event) error {
	return event.Update()
}

var DeleteEvent = func(event *models.Event) error {
	return event.Delete()
}
