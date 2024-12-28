package catalog

import (
	"context"
	"log"

	pb "github.com/Anshu-rai89/go-ms/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return &Client{conn: conn, service: c}, nil

}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {

	r, err := c.service.PostProduct(ctx, &pb.PostProductRequest{Name: name, Description: description, Price: price})

	if err != nil {
		return nil, err
	}

	return &Product{Name: r.Product.Name, Description: r.Product.Description, ID: r.Product.Id, Price: r.Product.Price}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{Id: id})

	if err != nil {
		return nil, err
	}

	return &Product{Name: r.Product.Name, Description: r.Product.Description, ID: r.Product.Id, Price: r.Product.Price}, nil

}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, p := range r.Product {
		products = append(products, &Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ID:          p.Id,
		})
	}

	return products, nil
}

func (c *Client) GetProductsByIds(ctx context.Context, ids []string) ([]*Product, error) {
	r, err := c.service.GetProductsByIds(ctx, &pb.GetProductsByIdsRequest{Ids: ids})

	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, p := range r.Product {
		products = append(products, &Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ID:          p.Id,
		})
	}

	return products, nil
}

func (c *Client) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	r, err := c.service.SearchProducts(ctx, &pb.SearchProductsRequest{Query: query, Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, p := range r.Product {
		products = append(products, &Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ID:          p.Id,
		})
	}

	return products, nil
}
