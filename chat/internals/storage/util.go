package storage

import (
	"slices"
	"webchat-server/chat/internals/domain"
)

func findDialogID(storage *Storage, login1 domain.LoginType, login2 domain.LoginType) (domain.DialogID, bool) {
	userDialogIDs, _ := storage.usersToDialogs[login1]
	dialogIdx := slices.IndexFunc(userDialogIDs, func(id domain.DialogID) bool {
		dialog, _ := storage.dialogs[id]
		return dialog.Contains(login2)
	})
	if dialogIdx == -1 {
		return domain.DialogID{}, false
	}
	return userDialogIDs[dialogIdx], true
}

func (storage *Storage) CreateDialog(login1 domain.LoginType, login2 domain.LoginType) domain.Dialog {
	dialog := domain.NewDialog(login1, login2)
	storage.dialogs[dialog.ID] = dialog
	firstUserDialogs, _ := storage.usersToDialogs[login1]
	firstUserDialogs = append(firstUserDialogs, dialog.ID)
	secondUserDialogs, _ := storage.usersToDialogs[login2]
	secondUserDialogs = append(secondUserDialogs, dialog.ID)
	return dialog
}
