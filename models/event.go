package models

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   string    `json:"startTime"`
	EndTime     string    `json:"endTime"`
	LocationID  int       `json:"locationId"`
	Location    *Location `json:"location"`
}

type EventInventoryEntry struct {
	ID      int `json:"id"`
	EventID int `json:"eventId"`
}
