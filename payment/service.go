package payment

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostPayment(ctx context.Context, orderId string, status string, amount float64) (*Payment, error)
	GetPayments(ctx context.Context, skip, take uint64) ([]*Payment, error)
}

type Payment struct {
	ID      string  `json:"id"`
	OrderID string  `json:"orderId"`
	Status  string  `json:"status"`
	Amount  float64 `json:"amount"`
}

type paymentService struct {
	repository Repository
}

func NewPaymentService(r Repository) Service {
	return &paymentService{repository: r}
}

func (s *paymentService) PostPayment(ctx context.Context, orderId string, status string, amount float64) (*Payment, error) {
	p := &Payment{
		ID:      ksuid.New().String(),
		OrderID: orderId,
		Status:  status,
		Amount:  amount,
	}

	err := s.repository.PutPayment(ctx, p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *paymentService) GetPayments(ctx context.Context, skip, take uint64) ([]*Payment, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}

	return s.repository.ListPayments(ctx, skip, take)
}
