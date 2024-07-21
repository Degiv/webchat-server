package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"webchat-server/chat/internals/domain"
)

var (
	errNoSuchUser   = errors.New("no such user")
	errNoSuchDialog = errors.New("no such dialog")
)

type Storage struct {
	database       *sql.DB
	logger         slog.Logger
	users          map[domain.LoginType]domain.User
	allChat        []byte
	dialogs        map[domain.DialogID]domain.Dialog
	usersToDialogs map[domain.LoginType][]domain.DialogID
	storageMutex   sync.RWMutex
}

func NewChatStorage(database *sql.DB, logger slog.Logger) *Storage {
	return &Storage{
		database:       database,
		logger:         logger,
		users:          make(map[domain.LoginType]domain.User),
		allChat:        make([]byte, 0),
		dialogs:        make(map[domain.DialogID]domain.Dialog),
		usersToDialogs: make(map[domain.LoginType][]domain.DialogID),
	}
}

func (storage *Storage) GetAllChat() []domain.Message {
	allChat := make([]domain.Message, 0)
	json.Unmarshal(storage.allChat, &allChat)
	return allChat
}

func (storage *Storage) CreateUser(login domain.LoginType) {
	storage.users[login] = domain.User{Login: login}
	storage.database.Exec("INSERT INTO users (login) VALUES (?)", login)
}

func (storage *Storage) UserExist(login domain.LoginType) bool {
	row := storage.database.QueryRow("SELECT * FROM users WHERE login = ?", login)
	return row.Scan() == nil
}

func (storage *Storage) PostAllChat(message domain.Message) {
	marshalledMessage, _ := json.Marshal(message)
	storage.allChat = append(storage.allChat, marshalledMessage...)
}

func (storage *Storage) GetPrivateChat(userLogin domain.LoginType, companionLogin domain.LoginType) []domain.Message {
	dialogID, ok := findDialogID(storage, userLogin, companionLogin)
	if !ok {
		return []domain.Message{}
	}
	return storage.dialogs[dialogID].Messages
}

func (storage *Storage) PostPrivateChat(senderLogin domain.LoginType, receiverLogin domain.LoginType, message domain.Message) {
	dialogID, ok := findDialogID(storage, senderLogin, receiverLogin)
	var dialog domain.Dialog
	if !ok {
		dialog = storage.CreateDialog(senderLogin, receiverLogin)
	} else {
		dialog, _ = storage.dialogs[dialogID]
	}
	dialog.AddMessage(message)
}
