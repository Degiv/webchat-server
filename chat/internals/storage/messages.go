package storage

import (
	"github.com/jmoiron/sqlx"
	"webchat-server/chat/internals/domain"
)

type Messages interface {
	Insert(message *domain.Message)
}

func NewMessages(DB *sqlx.DB) Messages {
	return &MessagesDB{DB: DB}
}

type MessagesDB struct {
	DB *sqlx.DB
}

func (database *MessagesDB) Insert(message *domain.Message) {
	const query = "INSERT INTO messages VALUES (?, ?, ?, ?)"
	database.DB.Exec(query,
		message.Id,
		message.UserLogin,
		message.ConversationId,
		message.Text)
}
