package payment

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	PutPayment(ctx context.Context, p *Payment) error
	ListPayments(ctx context.Context, skip, take uint64) ([]*Payment, error)
}

type postGresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgresql", url)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &postGresRepository{
		db: db,
	}, nil
}

func (r *postGresRepository) Close() {
	r.db.Close()
}

func (r *postGresRepository) PutPayment(ctx context.Context, p *Payment) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO payments(id, order_id, status, amount) VALUES($1, $2,$3, $4)",
		p.ID,
		p.OrderID,
		p.Status,
		p.Amount,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *postGresRepository) ListPayments(ctx context.Context, skip, take uint64) ([]*Payment, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, order_id,status, amount from payments offset $1 limit $2", skip, take)

	if err != nil {
		return nil, err
	}

	payments := []*Payment{}

	for rows.Next() {
		p := &Payment{}

		err = rows.Scan(&p.ID, &p.OrderID, &p.Status, &p.Amount)

		if err != nil {
			return nil, err
		}

		payments = append(payments, p)
	}

	return payments, nil
}
