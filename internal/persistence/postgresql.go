package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Next func(row pgx.Rows) error
type Transaction func (tx pgx.Tx) error

type Postgresql struct {
	Pool      *pgxpool.Pool
	Connected bool
	Url       string
}

func (p *Postgresql) Create() error {
	pool, err := pgxpool.New(context.Background(), p.Url)
	if err != nil {
		return err
	}
	p.Pool = pool
	p.Connected = true
	return nil
}

func (p *Postgresql) Query(next Next, query string, values ...any) error {
	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), query, values...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if next(rows) != nil {
			break
		}
	}
	return nil
}

func (p *Postgresql) Exec(query string, values ...any) (int64, error) {
	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	tag, err := conn.Exec(context.Background(), query, values...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (p *Postgresql) Txn(tx Transaction) error{
	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ptx, err := conn.BeginTx(context.Background(),pgx.TxOptions{})
	if err != nil{
		return err
	}
	err = tx(ptx)
	defer func() {
        if err != nil {
            ptx.Rollback(context.Background())
        } else {
            ptx.Commit(context.Background())
        }
    }()
	
	return nil
}

func (p *Postgresql) Close() {
	if !p.Connected {
		return
	}
	p.Pool.Close()
}
