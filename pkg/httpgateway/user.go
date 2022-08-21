package httpgateway

import (
	"coinChat/pkg/repository"
	"context"
	"fmt"
)

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Password  string `json:"password"`
}

func (user *User) ToRepoUser() *repository.User {
	return &repository.User{
		Id:        user.Id,
		Name:      user.Name,
		PublicKey: user.PublicKey,
		Password:  user.Password,
	}
}

func FromRepoUser(repoUser repository.User) *User {
	return &User{
		Id:        repoUser.Id,
		Name:      repoUser.Name,
		PublicKey: repoUser.PublicKey,
		Password:  repoUser.Password,
	}
}

func (h *Handler) SignUp(ctx context.Context, user User) (*User, error) {
	// todo: verify the user
	_, err := h.userRepo.GetUserByNameOrPublicKey(ctx, user.Name, user.PublicKey)
	if err == nil {
		return nil, fmt.Errorf("user already signed up or name is already taken")
	}
	if err := h.userRepo.AddUser(ctx, *user.ToRepoUser()); err != nil {
		return nil, err
	}
	userRepo, err := h.userRepo.GetUserByNameAndPassword(ctx, user.Name, user.Password)
	if err != nil {
		return nil, err
	}
	return FromRepoUser(*userRepo), nil
}

func (h *Handler) Login(ctx context.Context, user User) (*User, error) {
	dbUser, err := h.userRepo.GetUserByNameAndPassword(ctx, user.Name, user.Password)
	if err != nil || dbUser == nil {
		return nil, fmt.Errorf("invalid name or password")
	}
	return FromRepoUser(*dbUser), nil
}

func (h *Handler) Search(ctx context.Context, user User) (*User, error) {
	dbUser, err := h.userRepo.GetByPublicKey(ctx, user.PublicKey)
	if err != nil || dbUser == nil {
		return nil, fmt.Errorf("user is not found")
	}
	return FromRepoUser(*dbUser), nil
}
