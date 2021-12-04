package event

import (
	"bookstop/graph/model"
	"strconv"
)

func ToGraph(evt *Event) *model.Event {
	if evt == nil {
		return nil
	}

	return &model.Event{
		ID:          strconv.Itoa(int(evt.ID.Int)),
		Title:       evt.Title.String,
		Description: evt.Description.String,
		Href:        evt.Href.String,
		UserID:      strconv.Itoa(int(evt.UserID.Int)),
		StartedAt:   evt.StartedAt.Time,
		EndedAt:     evt.EndedAt.Time,
	}
}
