package handlers

import "webchat-server/chat/internals/domain"

type PostPrivateChatRequest struct {
	receiverLogin domain.LoginType
	message       domain.Message
}
