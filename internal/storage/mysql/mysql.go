package mysql

import (
	"context"
	"database/sql"
	"github.com/fredmayer/go-rest-api-template/pkg/database/mysql"
)

type Storage struct {
	db    *sql.DB
	Dummy *DummyRepository
}

func New(host string, port string, user string, password string, db string) *Storage {
	dbc := mysql.NewClient(context.Background(), mysql.Options{
		host, port, user, password, db,
	})

	return &Storage{
		db:    dbc,
		Dummy: NewDummy(dbc),
	}
}

func (s *Storage) Stop() {
	s.db.Close()
}
