package levelDB

import (
	"context"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type levelDb struct {
	db *leveldb.DB
}

type DatabaseI interface {
	Db(ctx context.Context) interface{}
}

func NewConnection(dsn string) (DatabaseI, error) {
	var err error

	db, err := leveldb.OpenFile(dsn, nil)
	if err != nil {
		return levelDb{db: db}, err
	}

	log.Printf("successfully connected to database: %v", dsn)

	return &levelDb{db: db}, nil
}

func (p levelDb) Db(ctx context.Context) interface{} {
	tx := ctx.Value("txContext")
	if tx == nil {
		return p.db
	}
	return tx.(*leveldb.DB)
}
