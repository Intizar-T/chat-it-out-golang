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
	CreatedAt    time.Time `json:"createdAt"`
}

type Conversation struct {
	UserId   int `json:"userId"`
	FriendId int `json:"friendId"`
}

func MessageFromRepoMessage(repoMessage *repository.Message) Message {
	return Message{
		Id:           repoMessage.Id,
		SenderId:     repoMessage.SenderId,
		SenderName:   repoMessage.SenderName,
		ReceiverId:   repoMessage.ReceiverId,
		ReceiverName: repoMessage.ReceiverName,
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

func (h *Handler) GetMessages(ctx context.Context, user User) (map[int][]Message, error) {
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
		if _, ok := res[m.SenderId]; ok {
			res[friendId] = []Message{}
		}
		res[m.SenderId] = append(res[m.SenderId], MessageFromRepoMessage(&m))
	}
	return res, nil
}

func (h *Handler) NewMessage(ctx context.Context, message Message) error {
	return h.messageRepo.CreateMessage(ctx, message.ConvertToMessageRepo())
}
