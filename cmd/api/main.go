package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"webchat-server/chat/internals/handlers"
	"webchat-server/chat/internals/services"
	"webchat-server/chat/internals/storage"
)

func main() {
	logger := slog.Logger{}
	dsn := "postgres://postgres:postgres@localhost:5432/WebchatDB" + "?sslmode=disable"
	database, _ := sql.Open("postgres", dsn)

	chatStorage := storage.NewChatStorage(database, logger)

	chatService := services.NewChatService(logger, chatStorage)

	chatHandler := handlers.NewChatHandler(logger, chatService)

	r := chi.NewRouter()

	r.Mount("/chat", chatHandler.Routes())

	http.ListenAndServe(":3000", r)
}
