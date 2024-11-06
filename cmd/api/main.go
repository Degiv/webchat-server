package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"webchat-server/chat/internals/handlers"
	"webchat-server/chat/internals/services"
	"webchat-server/chat/internals/storage"
)

func main() {
	//logger := slog.Logger{}

	db, err := sqlx.Connect("postgres", "user=postgres dbname=WebchatDB sslmode=disable")

	if err != nil {
		log.Fatalln(err)
	}

	db.Ping()
	MessagesDB := storage.NewMessages(db)

	chatStorage := storage.NewChatStorage(database, logger)

	chatService := services.NewChatService(logger, chatStorage)

	chatHandler := handlers.NewChatHandler(logger, chatService)

	r := chi.NewRouter()

	r.Mount("/chat", chatHandler.Routes())

	http.ListenAndServe(":3000", r)
}
