//go:generate protoc ./account.proto --go_out=plugins=grpc:./pb

package account

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Anshu-rai89/go-ms/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Println("Error in ListenGRPC", err)
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)

	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{Account: &pb.Account{
		Id:   a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)

	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{Account: &pb.Account{
		Id:   a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	resp, err := s.service.GetAccounts(ctx, r.Skip, r.Take)

	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, p := range resp {
		accounts = append(accounts, &pb.Account{
			Name: p.Name,
			Id:   p.ID,
		})
	}

	return &pb.GetAccountsResponse{
		Account: accounts,
	}, nil
}
