package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Anshu-rai89/go-ms/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Println("Error in catalog ListenGRPC", err)

		return err
	}

	serv := grpc.NewServer()
	pb.RegisterProductServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)

	if err != nil {
		return nil, err
	}

	product := &pb.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Id:          p.ID,
	}

	return &pb.PostProductResponse{
		Product: product,
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)

	if err != nil {
		return nil, err
	}

	product := &pb.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Id:          p.ID,
	}

	return &pb.GetProductResponse{
		Product: product,
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	resp, err := s.service.GetProducts(ctx, r.Skip, r.Take)

	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}

	for _, p := range resp {
		products = append(products, &pb.Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Id:          p.ID,
		})
	}

	return &pb.GetProductsResponse{
		Product: products,
	}, nil
}

func (s *grpcServer) GetProductsByIds(ctx context.Context, r *pb.GetProductsByIdsRequest) (*pb.GetProductsByIdsResponse, error) {
	resp, err := s.service.GetProductsByIds(ctx, r.Ids)

	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}

	for _, p := range resp {
		products = append(products, &pb.Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Id:          p.ID,
		})
	}

	return &pb.GetProductsByIdsResponse{
		Product: products,
	}, nil
}

func (s *grpcServer) SearchProducts(ctx context.Context, r *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	resp, err := s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)

	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}

	for _, p := range resp {
		products = append(products, &pb.Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Id:          p.ID,
		})
	}

	return &pb.SearchProductsResponse{
		Product: products,
	}, nil
}
