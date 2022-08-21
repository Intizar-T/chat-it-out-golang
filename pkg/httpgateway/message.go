package httpgateway

import (
	"coinChat/pkg/repository"
	"context"
	"time"
)

type Message struct {
	Id           int       `json:"id"`
	SenderId     int       `json:"senderId"`
	SenderName   string    `json:"senderName"`
	ReceiverId   int       `json:"receiverId"`
	ReceiverName string    `json:"receiverName"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Conversation struct {
	UserId   int `json:"userId"`
	FriendId int `json:"friendId"`
}

type ConversationMessages struct {
	Id       int       `json: "id"`
	Messages []Message `json:"messages"`
}

func MessageFromRepoMessage(repoMessage *repository.Message) Message {
	return Message{
		Id:           repoMessage.Id,
		SenderId:     repoMessage.SenderId,
		SenderName:   repoMessage.SenderName,
		ReceiverId:   repoMessage.ReceiverId,
		ReceiverName: repoMessage.ReceiverName,
		Message:      repoMessage.Message,
		CreatedAt:    repoMessage.CreatedAt,
	}
}

func (m *Message) ConvertToMessageRepo() repository.Message {
	return repository.Message{
		Id:           m.Id,
		SenderId:     m.SenderId,
		SenderName:   m.SenderName,
		ReceiverId:   m.ReceiverId,
		ReceiverName: m.ReceiverName,
		Message:      m.Message,
		CreatedAt:    m.CreatedAt,
	}
}

func (h *Handler) GetConversations(ctx context.Context, userId int, friendId int) ([]Message, error) {
	mes, err := h.messageRepo.GetConversations(ctx, userId, friendId)
	if err != nil {
		return nil, err
	}
	res := []Message{}
	for _, m := range mes {
		res = append(res, MessageFromRepoMessage(&m))
	}
	return res, nil
}

func (h *Handler) GetMessages(ctx context.Context, user User) ([]ConversationMessages, error) {
	mes, err := h.messageRepo.GetSenderMessages(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	res := make(map[int][]Message)
	for _, m := range mes {
		friendId := m.SenderId
		if friendId == user.Id {
			friendId = m.ReceiverId
		}
		if _, ok := res[friendId]; !ok {
			res[friendId] = []Message{}
		}
		res[friendId] = append(res[friendId], MessageFromRepoMessage(&m))
	}
	cm := []ConversationMessages{}
	i := 0
	for key, val := range res {
		cm = append(cm, ConversationMessages{
			Id:       key,
			Messages: []Message{},
		})
		for _, mes := range val {
			cm[i].Messages = append(cm[i].Messages, mes)
		}
		i += 1
	}
	return cm, nil
}

func (h *Handler) NewMessage(ctx context.Context, message Message) error {
	return h.messageRepo.CreateMessage(ctx, message.ConvertToMessageRepo())
}
