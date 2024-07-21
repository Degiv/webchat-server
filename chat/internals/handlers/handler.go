package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"webchat-server/chat/internals/domain"
)

const (
	cookieAuth = "auth"
	login      = "login"
)

type chatService interface {
	UserExist(login domain.LoginType) bool
	CreateUser(login domain.LoginType)

	GetAllChat() []domain.Message
	PostAllChat(message domain.Message)

	GetPrivateChat(userLogin domain.LoginType, companionLogin domain.LoginType) []domain.Message
	PostPrivateChat(senderLogin domain.LoginType, receiverLogin domain.LoginType, message domain.Message)
}

type ChatHandler struct {
	logger      slog.Logger
	chatService chatService
}

func NewChatHandler(logger slog.Logger, chatService chatService) *ChatHandler {
	return &ChatHandler{
		logger:      logger,
		chatService: chatService,
	}
}

func (handler *ChatHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", handler.Login)

	r.Get("/all", handler.GetAllChat)

	userRouter := chi.NewRouter()
	r.Mount("/user", userRouter)
	userRouter.Use(Auth)
	userRouter.Post("/all", handler.PostAllChat)
	userRouter.Get("/message", handler.GetPrivateChat)
	userRouter.Post("/message", handler.PostPrivateChat)

	return r
}

func (handler *ChatHandler) GetAllChat(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(handler.chatService.GetAllChat())
	w.Write(data)
}

func (handler *ChatHandler) PostAllChat(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	message := domain.Message{}
	json.Unmarshal(data, &message)
	handler.chatService.PostAllChat(message)
}

func (handler *ChatHandler) GetPrivateChat(w http.ResponseWriter, r *http.Request) {
	bodyData, _ := io.ReadAll(r.Body)
	companionLogin := domain.LoginType("")
	json.Unmarshal(bodyData, &companionLogin)
	userLogin := r.Context().Value(login).(domain.LoginType)

	dialogData, _ := json.Marshal(handler.chatService.GetPrivateChat(userLogin, companionLogin))
	w.Write(dialogData)
}

func (handler *ChatHandler) PostPrivateChat(w http.ResponseWriter, r *http.Request) {
	bodyData, _ := io.ReadAll(r.Body)
	unmarshalled := PostPrivateChatRequest{}
	json.Unmarshal(bodyData, &unmarshalled)
	senderLogin := r.Context().Value(login).(domain.LoginType)
	handler.chatService.PostPrivateChat(senderLogin, unmarshalled.receiverLogin, unmarshalled.message)
}

func Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieAuth)
		switch err {
		case nil:
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if c.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idCtx := context.WithValue(r.Context(), login, domain.LoginType(c.Value))

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
}

func (handler *ChatHandler) Login(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header()
		return
	}
	defer r.Body.Close()

	var u domain.User
	err = json.Unmarshal(d, &u)
	if err != nil || u.Login == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := &http.Cookie{
		Name:  cookieAuth,
		Value: string(u.Login),
		Path:  "/",
	}

	if !handler.chatService.UserExist(u.Login) {
		handler.chatService.CreateUser(u.Login)
	}

	http.SetCookie(w, c)
}
