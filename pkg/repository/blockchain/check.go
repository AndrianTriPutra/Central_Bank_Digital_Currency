package blockchain

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
)

func (bc blockchain) isGenesis(ctx context.Context, tx *domain.Transaction) bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].Out == -1
}

func (bc blockchain) canUnlock(ctx context.Context, in domain.TxInput, data string) bool {
	return in.Sig == data
}

func (bc blockchain) canBeUnlocked(ctx context.Context, out domain.TxOutput, data string) bool {
	return out.PubKey == data
}
