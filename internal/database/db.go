package database

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/samarthasthan/tanx-task/internal/database/mysql/sqlc"
)

type Database interface {
	Connect(string) error
	Close() error
}

type MySQL struct {
	Queries *sqlc.Queries
	DB      *sql.DB
}

type Redis struct {
	*redis.Client
}

func NewMySQL() Database {
	return &MySQL{}
}

func NewRedis() *Redis {
	return &Redis{}
}

func (s *MySQL) Connect(addr string) error {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return err
	}
	s.DB = db
	s.Queries = sqlc.New(db)
	return nil
}

func (s *MySQL) Close() error {
	return s.DB.Close()
}

func (r *Redis) Connect(addr string) error {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	r.Client = rdb
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Close() error {
	err := r.Close()
	if err != nil {
		return err
	}
	return nil
}
