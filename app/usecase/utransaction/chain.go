package utransaction

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
)

func (u usecase) GetChain(ctx context.Context) ([]domain.Block, error) {
	data, err := u.repoCrud.GetChain(ctx)
	if err != nil {
		return data, err
	}

	return data, nil
}
