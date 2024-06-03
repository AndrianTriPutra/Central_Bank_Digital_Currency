package crud

import (
	"atp/cbdc/pkg/adapter/levelDB"
	"atp/cbdc/pkg/utils/domain"
	"context"
	"encoding/json"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

type repository struct {
	provider levelDB.DatabaseI
}

func NewRepository(provider levelDB.DatabaseI) RepositoryI {
	return repository{
		provider,
	}
}

type RepositoryI interface {
	PUT(ctx context.Context, key, value string) error
	GET(ctx context.Context, key string) (string, error)
	GetChain(ctx context.Context) ([]domain.Block, error)
	GetWallet(ctx context.Context) ([]string, error)
}

func (r repository) PUT(ctx context.Context, key, value string) error {
	db := r.provider.Db(ctx).(*leveldb.DB)
	err := db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) GET(ctx context.Context, key string) (string, error) {
	var value string
	db := r.provider.Db(ctx).(*leveldb.DB)
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return value, err
	}
	value = string(data)
	return value, nil
}

func (r repository) GetChain(ctx context.Context) ([]domain.Block, error) {
	var data []domain.Block
	value, err := r.GET(ctx, "lh")
	if err != nil {
		return data, err
	}

	lastkey := value
	for {
		value, err = r.GET(ctx, lastkey)
		if err != nil {
			buff := err.Error()
			if strings.Contains(buff, "not found") {
				break
			} else {
				return data, err
			}
		}

		if strings.Contains(value, "transaction") {
			var block domain.Block
			err = json.Unmarshal([]byte(value), &block)
			if err != nil {
				return data, err
			}
			data = append(data, block)
			lastkey = block.Header.PrevHash
			if block.Header.PrevHash == "0000000000000000000000000000000000000000000000000000000000000000" {
				return data, nil
			}

		}
	}
	return data, nil
}

func (r repository) GetWallet(ctx context.Context) ([]string, error) {
	var wallet domain.DWallet

	data, err := r.GET(ctx, "wallet")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &wallet)
	if err != nil {
		return nil, err
	}

	return wallet.Address, nil
}
