package catalog

import (
	"context"
	"encoding/json"
	"errors"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	errorNotFound = errors.New("Product not found")
)

type Repository interface {
	Close()
	GetProductById(ctx context.Context, id string) (*Product, error)
	PutProduct(ctx context.Context, p Product) error
	ListProduct(ctx context.Context, skip, take uint64) ([]*Product, error)
	ListProductsWithIDS(ctx context.Context, ids []string) ([]*Product, error)
	SearchProduct(ctx context.Context, query string, skip, take uint64) ([]*Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)

	if err != nil {
		return nil, err
	}

	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {

}

func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().Index("catalog").Type("product").Id(id).Do(ctx)

	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, errorNotFound
	}
	p := productDocument{}
	if err = json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
	}, nil

}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().Index("catalog").Type("product").Id(p.ID).BodyJson(productDocument{
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
	}).Do(ctx)

	return err
}

func (r *elasticRepository) ListProduct(ctx context.Context, skip, take uint64) ([]*Product, error) {
	resp, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMatchAllQuery()).
		From(int(skip)).Size(int(take)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, hit := range resp.Hits.Hits {
		p := &productDocument{}

		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, &Product{
				ID:          hit.Id,
				Name:        p.Name,
				Price:       p.Price,
				Description: p.Description,
			})
		}

		if err != nil {
			return nil, err
		}

	}

	return products, nil
}

func (r *elasticRepository) ListProductsWithIDS(ctx context.Context, ids []string) ([]*Product, error) {
	items := []*elastic.MultiGetItem{}

	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().
			Index("catalog").
			Type("product").
			Id(id))
	}

	resp, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, doc := range resp.Docs {
		p := &productDocument{}

		if err = json.Unmarshal(*doc.Source, &p); err == nil {
			products = append(products, &Product{
				ID:          doc.Id,
				Name:        p.Name,
				Price:       p.Price,
				Description: p.Description,
			})
		}

		if err != nil {
			return nil, err
		}

	}
	return products, nil
}

func (r *elasticRepository) SearchProduct(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	resp, err := r.client.Search().Index("catalog").Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).From(int(skip)).Size(int(take)).Do(ctx)

	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, hit := range resp.Hits.Hits {
		p := &productDocument{}

		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, &Product{
				ID:          hit.Id,
				Name:        p.Name,
				Price:       p.Price,
				Description: p.Description,
			})
		}

		if err != nil {
			return nil, err
		}

	}

	return products, nil
}
