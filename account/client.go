package account

import (
	"context"
	"log"

	pb "github.com/Anshu-rai89/go-ms/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c := pb.NewAccountServiceClient(conn)

	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})

	if err != nil {
		return nil, err
	}

	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	accounts := []Account{}

	for _, a := range r.Account {
		accounts = append(accounts, Account{ID: a.Id, Name: a.Name})
	}

	return accounts, nil
}