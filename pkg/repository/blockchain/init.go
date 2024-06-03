package blockchain

import (
	"atp/cbdc/pkg/repository/crud"
	"atp/cbdc/pkg/repository/pow"
	"atp/cbdc/pkg/utils/domain"
	"context"
)

type Genesis struct {
	Address string
	Balance int
}

type blockchain struct {
	setting  Genesis
	repoCrud crud.RepositoryI
	repoPow  pow.RepositoryI
}

func NewRepository(setting Genesis, repoCrud crud.RepositoryI, repoPow pow.RepositoryI) RepositoryI {
	return &blockchain{
		setting:  setting,
		repoCrud: repoCrud,
		repoPow:  repoPow,
	}
}

type RepositoryI interface {
	InitBlockchain(ctx context.Context) error
	CreateBlock(ctx context.Context, a *domain.Block) error

	FindUnspentTransactions(ctx context.Context, address string) ([]domain.Transaction, error)
	FindUTXO(ctx context.Context, address string) ([]domain.TxOutput, error)

	Send(ctx context.Context, trans domain.Trans) error
	NewTransaction(ctx context.Context, trans domain.Trans) (domain.Transaction, error)
	FindSpendableOutputs(ctx context.Context, trans domain.Trans) (int, map[string][]int, error)
}
