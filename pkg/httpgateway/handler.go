package httpgateway

import (
	"coinChat/pkg/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	userRepo    *repository.UserRepo
	messageRepo *repository.MessageRepo
}

func NewHandler(userRepo *repository.UserRepo, messageRepo *repository.MessageRepo) Handler {
	return Handler{
		userRepo:    userRepo,
		messageRepo: messageRepo,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.HandleGet(ctx, w, r)
	case http.MethodPost:
		h.HandlePost(ctx, w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.WriteMessage(w, fmt.Sprintf("Method is not defined: %+v", r.Method))
	}
}

func unMarshall(body io.ReadCloser, val interface{}) error {
	b, _ := ioutil.ReadAll(body)
	return json.Unmarshal([]byte(b), val)
}

func (h *Handler) WriteMessage(w http.ResponseWriter, res interface{}) {
	resp := make(map[string]interface{})
	resp["message"] = res
	b, _ := json.Marshal(resp)
	io.WriteString(w, string(b))
}

func (h *Handler) Write(w http.ResponseWriter, res interface{}) {
	b, _ := json.Marshal(res)
	io.WriteString(w, string(b))
}

func (h *Handler) HandlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		var user User
		if err := unMarshall(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		loggedUser, err := h.Login(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.Write(w, err.Error())
			return
		}
		h.Write(w, loggedUser)
	case "/signup":
		var user User
		if err := unMarshall(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		signedUser, err := h.SignUp(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		h.Write(w, signedUser)
		return
	case "/newmessage":
		var message Message
		if err := unMarshall(r.Body, &message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		err := h.NewMessage(ctx, message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		h.WriteMessage(w, "OK")
		return
	default:
		// w.WriteHeader(http.StatusNotFound)
		h.WriteMessage(w, "status not found")
		return
	}
}

func (h *Handler) HandleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/messages":
		var user User
		if err := unMarshall(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		messages, err := h.GetMessages(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err)
			return
		}
		h.WriteMessage(w, messages)
	case "/conversations":
		var conversation Conversation
		if err := unMarshall(r.Body, &conversation); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		messages, err := h.GetConversations(ctx, conversation.UserId, conversation.FriendId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err)
			return
		}
		h.WriteMessage(w, messages)
	default:
		// w.WriteHeader(http.StatusNotFound)
		h.WriteMessage(w, "invalid path for get request")
		return
	}
}
