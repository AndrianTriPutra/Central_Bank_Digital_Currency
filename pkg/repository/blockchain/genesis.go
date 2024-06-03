package blockchain

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
	"fmt"
)

func (bc blockchain) genesis(ctx context.Context) (domain.Transaction, error) {
	var inputs []domain.TxInput
	input := domain.TxInput{
		ID:  fmt.Sprintf("%x", [32]byte{}),
		Out: -1,
		Sig: "first devisa from Underlying Gold",
	}
	inputs = append(inputs, input)

	var outputs []domain.TxOutput
	output := domain.TxOutput{
		PubKey: bc.setting.Address,
		Value:  bc.setting.Balance,
	}
	outputs = append(outputs, output)

	trans := domain.Transaction{
		ID:      "",
		Inputs:  inputs,
		Outputs: outputs,
	}

	txID := trans.Hash()
	trans.ID = txID

	return trans, nil
}
