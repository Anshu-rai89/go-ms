package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/Anshu-rai89/go-ms/account"
	"github.com/Anshu-rai89/go-ms/catalog"
	"github.com/Anshu-rai89/go-ms/order"
	"github.com/Anshu-rai89/go-ms/payment"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
	paymentClient *payment.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl, paymentUrl string) (*Server, error) {
	accountClient, err := account.NewClient(accountUrl)

	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		orderClient.Close()
		return nil, err
	}

	paymentClient, err := payment.NewClient(paymentUrl)
	if err != nil {
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
		paymentClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}
func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
