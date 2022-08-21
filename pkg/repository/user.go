package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id         int       `db:"id"`
	Name       string    `db:"user_name"`
	PublicKey  string    `db:"public_key"`
	Password   string    `db:"user_password"`
	Created_at time.Time `db:"created_at"`
}

type UserRepo struct {
	BaseRepo
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		&Repository{
			DB: db,
		},
	}
}

func (r *UserRepo) AddUser(ctx context.Context, user User) error {
	query := `
	INSERT INTO users
	(user_name, public_key, user_password) VALUES
	(:user_name, :public_key, :user_password)
	`
	return r.CreateOne(ctx, query, user)
}

func (r *UserRepo) GetUser(ctx context.Context, userId int64) (*User, error) {
	query := "select * from users where id=$1"
	var result User
	err := r.QueryOne(ctx, &result, query, userId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserRepo) GetByPublicKey(ctx context.Context, publicKey string) (*User, error) {
	query := "select * from users where public_key=$1"
	var result User
	err := r.QueryOne(ctx, &result, query, publicKey)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserRepo) GetUserByNameAndPassword(ctx context.Context, name string, pwd string) (*User, error) {
	query := "select * from users where user_name=$1 AND user_password=$2"
	var result User
	err := r.QueryOne(ctx, &result, query, name, pwd)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserRepo) GetUserByNameOrPublicKey(ctx context.Context, name string, publicKey string) (*User, error) {
	query := "select * from users where user_name=$1 OR public_key=$2"
	var result User
	err := r.QueryOne(ctx, &result, query, name, publicKey)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
