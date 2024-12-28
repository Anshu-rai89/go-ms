package payment

import (
	"context"

	pb "github.com/Anshu-rai89/go-ms/payment/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.PaymentServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	s := pb.NewPaymentServiceClient(conn)

	return &Client{conn: conn, service: s}, nil
}

func (c *Client) PostPayment(ctx context.Context, orderId string, status string, amount float64) (*Payment, error) {
	res, err := c.service.PostPayment(ctx, &pb.PostPaymentRequest{
		Status:  status,
		OrderId: orderId,
		Amount:  amount,
	})

	if err != nil {
		return nil, err
	}

	return &Payment{
		ID:      res.Payment.Id,
		OrderID: res.Payment.OrderId,
		Status:  res.Payment.Status,
		Amount:  res.Payment.Amount,
	}, nil
}

func (c *Client) GetPayments(ctx context.Context, skip uint64, take uint64) ([]*Payment, error) {
	resp, err := c.service.GetPayments(ctx, &pb.GetPaymentsRequest{Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	payments := []*Payment{}

	for _, p := range resp.Payment {
		p1 := &Payment{
			ID:      p.Id,
			OrderID: p.OrderId,
			Status:  p.Status,
			Amount:  p.Amount,
		}

		payments = append(payments, p1)
	}

	return payments, nil
}
