package services

import (
	"log/slog"
	"webchat-server/chat/internals/domain"
)

type chatService struct {
	logger  slog.Logger
	storage storageInterface
}

type storageInterface interface {
	CreateUser(login domain.LoginType)
	userExist(login domain.LoginType) bool

	GetAllChat() []domain.Message
	PostAllChat(message domain.Message)

	GetPrivateChat(userLogin domain.LoginType, companionLogin domain.LoginType) []domain.Message
	PostPrivateChat(senderLogin domain.LoginType, receiverLogin domain.LoginType, message domain.Message)
}

func NewChatService(logger slog.Logger, storage storageInterface) *chatService {
	return &chatService{
		logger:  logger,
		storage: storage,
	}
}

func (chatService *chatService) CreateUser(login domain.LoginType) {
	chatService.storage.CreateUser(login)
}

func (chatService *chatService) userExist(login domain.LoginType) bool {
	return chatService.storage.userExist(login)
}

func (chatService *chatService) GetAllChat() []domain.Message {
	return chatService.storage.GetAllChat()
}

func (chatService *chatService) PostAllChat(message domain.Message) {
	chatService.storage.PostAllChat(message)
}

func (chatService *chatService) GetPrivateChat(userLogin domain.LoginType, companionLogin domain.LoginType) []domain.Message {
	return chatService.storage.GetPrivateChat(userLogin, companionLogin)
}

func (chatService *chatService) PostPrivateChat(senderLogin domain.LoginType, receiverLogin domain.LoginType, message domain.Message) {
	chatService.storage.PostPrivateChat(senderLogin, receiverLogin, message)
}
