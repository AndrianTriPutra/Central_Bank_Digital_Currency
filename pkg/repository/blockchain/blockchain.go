package blockchain

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
)

func (bc blockchain) InitBlockchain(ctx context.Context) error {
	nb := new(domain.Block)
	trans, err := bc.genesis(ctx)
	if err != nil {
		return err
	}

	nb.GiveData(trans)
	err = bc.CreateBlock(ctx, nb)
	if err != nil {
		return err
	}

	return nil
}
