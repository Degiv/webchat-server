package main

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"webchat-server/chat/internals/handlers"
	"webchat-server/chat/internals/services"
	"webchat-server/chat/internals/storage"
)

func main() {
	logger := slog.Logger{}
	chatStorage := storage.NewChatStorage(logger)

	chatService := services.NewChatService(logger, chatStorage)

	chatHandler := handlers.NewChatHandler(logger, chatService)

	r := chi.NewRouter()

	r.Mount("/chat", chatHandler.Routes())

	http.ListenAndServe(":3000", r)
}
