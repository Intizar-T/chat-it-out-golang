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
		if message.ReceiverId == message.SenderId {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, "Sender and receiver are same")
		}
		err := h.NewMessage(ctx, message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		messages, err := h.GetConversations(ctx, message.ReceiverId, message.SenderId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err)
			return
		}
		h.WriteMessage(w, messages)
		return
	case "/messages":
		var b []byte
		r.Body.Read(b)
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
			h.Write(w, err)
			return
		}
		h.WriteMessage(w, messages)
	case "/search":
		var user User
		if err := unMarshall(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		signedUser, err := h.Search(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.WriteMessage(w, err.Error())
			return
		}
		h.Write(w, signedUser)
	default:
		// w.WriteHeader(http.StatusNotFound)
		h.WriteMessage(w, "status not found")
		return
	}
}
