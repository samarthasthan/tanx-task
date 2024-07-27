package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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

func NewMySQL() Database {
	return &MySQL{}
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
