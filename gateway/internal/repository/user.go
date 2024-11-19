package repository

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
)

type User struct {
	Login    string
	Password string
}

type UserRepository struct {
	users map[string]User
	mu    *sync.RWMutex
}

func NewUserRepository() UserRepository {
	return UserRepository{
		users: make(map[string]User),
		mu:    &sync.RWMutex{},
	}
}

func (repo *UserRepository) AddUser(user User) error {
	repo.mu.RLock()
	if _, exists := repo.users[user.Login]; exists {
		return ErrUserAlreadyExist
	}

	repo.mu.RUnlock()

	repo.mu.Lock()
	repo.users[user.Login] = user
	repo.mu.Unlock()

	return nil
}

func (repo *UserRepository) GetUser(_ context.Context, login string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[login]
	if !exists {
		return User{}, ErrUserNotFound
	}

	return user, nil
}
