package payment

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Anshu-rai89/go-ms/payment/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedPaymentServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	serv := grpc.NewServer()

	pb.RegisterPaymentServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostPayment(ctx context.Context, r *pb.PostPaymentRequest) (*pb.PostPaymentResponse, error) {
	resp, err := s.service.PostPayment(ctx, r.OrderId, r.Status, r.Amount)

	if err != nil {
		return nil, err
	}

	p := &pb.Payment{
		Id:      resp.ID,
		OrderId: resp.OrderID,
		Status:  resp.Status,
		Amount:  resp.Amount,
	}

	return &pb.PostPaymentResponse{
		Payment: p,
	}, nil
}

func (s *grpcServer) GetPayments(ctx context.Context, r *pb.GetPaymentsRequest) (*pb.GetPaymentsResponse, error) {
	resp, err := s.service.GetPayments(ctx, r.Skip, r.Take)

	if err != nil {
		return nil, err
	}

	payments := []*pb.Payment{}

	for _, p := range resp {
		p1 := &pb.Payment{
			Id:      p.ID,
			OrderId: p.OrderID,
			Status:  p.Status,
			Amount:  p.Amount,
		}

		payments = append(payments, p1)
	}

	return &pb.GetPaymentsResponse{Payment: payments}, nil
}
