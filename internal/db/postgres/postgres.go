package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn interface {
	Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

type Storage struct {
	conn Conn
}

func MustConnect(ctx context.Context, connStr string) *Storage {

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &Storage{
		conn: conn,
	}
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) SaveURL(ctx context.Context, url, alias string) error {

	query := "INSERT INTO urls(url, alias) VALUES ($1, $2)"
	_, err := s.conn.Query(ctx, query, url, alias)
	if err != nil {
		return fmt.Errorf("url save failed: %s", err.Error())
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {

	query := "SELECT url FROM urls WHERE alias=$1"

	var url string
	err := s.conn.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("url read failed: %s", err.Error())
	}

	return url, nil
}
