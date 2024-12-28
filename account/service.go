package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	GetAccount(ctx context.Context, id string) (*Account, error)
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &AccountService{repository: r}
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *AccountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}

	err := s.repository.PutAccount(ctx, *account)
	return account, err
}

func (s *AccountService) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
