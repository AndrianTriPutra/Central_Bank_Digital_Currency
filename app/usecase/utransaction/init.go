package utransaction

import (
	"atp/cbdc/pkg/repository/blockchain"
	"atp/cbdc/pkg/repository/crud"
	"atp/cbdc/pkg/utils/domain"
	"context"
)

type usecase struct {
	repoBC   blockchain.RepositoryI
	repoCrud crud.RepositoryI
}

func NewUsecase(repoBC blockchain.RepositoryI, repoCrud crud.RepositoryI) UsecaseI {
	return &usecase{
		repoBC:   repoBC,
		repoCrud: repoCrud,
	}
}

type UsecaseI interface {
	GetChain(ctx context.Context) ([]domain.Block, error)

	Transaction(ctx context.Context, trans domain.Trans) error
	GetBalance(ctx context.Context, addr string) (int, error)

	CheckWallet(ctx context.Context, addr string) ([]string, bool, error)
	GetWallet(ctx context.Context) ([]domain.Wallet, error)
	UpdateWallet(ctx context.Context, addr []string) error

	History(ctx context.Context, addr string) ([]domain.Block, error)
}
