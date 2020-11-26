package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage interface {
	GetNoticeByHash(h string)
}

type DBStorage struct {
	Db pgxpool.Pool
}

func NewStorage(d pgxpool.Pool) DBStorage {
	return DBStorage{Db: d}
}
