package domain

import "github.com/google/uuid"

type DialogID uuid.UUID

type Dialog struct {
	ID         DialogID
	Messages   []Message
	Companions map[LoginType]struct{}
}

func NewDialog(companionsSlice ...LoginType) Dialog {
	companions := make(map[LoginType]struct{})
	for _, companion := range companionsSlice {
		companions[companion] = struct{}{}
	}

	return Dialog{
		ID:         DialogID(uuid.New()),
		Messages:   make([]Message, 0),
		Companions: companions,
	}
}

func (dialog *Dialog) AddMessage(message Message) {
	dialog.Messages = append(dialog.Messages, message)
}

func (dialog *Dialog) Contains(login LoginType) bool {
	for l := range dialog.Companions {
		if l == login {
			return true
		}
	}
	return false
}
