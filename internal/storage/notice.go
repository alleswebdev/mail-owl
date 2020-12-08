package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBStorage struct {
	Db pgxpool.Pool
}

func NewStorage(d pgxpool.Pool) DBStorage {
	return DBStorage{Db: d}
}
