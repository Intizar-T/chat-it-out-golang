package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Message struct {
	Id           int       `db:"id"`
	SenderId     int       `db:"sender_id"`
	SenderName   string    `db:"sender_name"`
	ReceiverName string    `db:"receiver_name"`
	ReceiverId   int       `db:"receiver_id"`
	Message      string    `db:"user_message"`
	CreatedAt    time.Time `db:"created_at"`
}

type MessageRepo struct {
	BaseRepo
}

func NewMessageRepo(db *sqlx.DB) *MessageRepo {
	return &MessageRepo{
		&Repository{
			DB: db,
		},
	}
}

func (r *MessageRepo) CreateMessage(ctx context.Context, message Message) error {
	query := `
	INSERT INTO messages
	(sender_id, sender_name, receiver_id, receiver_name, user_message) VALUES
	(:sender_id, :sender_name, :receiver_id, :receiver_name, :user_message)
	`
	return r.CreateOne(ctx, query, message)
}

func (r *MessageRepo) GetSenderMessages(ctx context.Context, id int) ([]Message, error) {
	query := `
	SELECT * FROM messages
	WHERE receiver_id=$1 OR sender_id=$1 ORDER BY created_at ASC
	`
	var result []Message
	err := r.Query(ctx, &result, query, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *MessageRepo) GetConversations(ctx context.Context, userId, friendId int) ([]Message, error) {
	query := `
	SELECT * FROM messages
	WHERE (receiver_id=$1 AND sender_id=$2) OR (sender_id=$1 AND receiver_id=$2) ORDER BY created_at ASC
	`
	var result []Message
	err := r.Query(ctx, &result, query, userId, friendId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
