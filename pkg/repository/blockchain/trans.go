package blockchain

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
	"errors"
)

func (bc blockchain) Send(ctx context.Context, trans domain.Trans) error {
	tranS, err := bc.NewTransaction(ctx, trans)
	if err != nil {
		return err
	}

	nb := new(domain.Block)
	nb.GiveData(tranS)
	err = bc.CreateBlock(ctx, nb)
	if err != nil {
		return err
	}
	return nil
}

func (bc blockchain) NewTransaction(ctx context.Context, trans domain.Trans) (domain.Transaction, error) {
	var inputs []domain.TxInput
	var outputs []domain.TxOutput
	var tx domain.Transaction

	acc, validOutputs, err := bc.FindSpendableOutputs(ctx, trans)
	if err != nil {
		return tx, err
	}

	if acc < trans.Amount {
		err = errors.New("acc not enough funds")
		return tx, err
	}

	for txid, outs := range validOutputs {
		for _, out := range outs {
			input := domain.TxInput{
				ID:  txid,
				Out: out,
				Sig: trans.From}
			inputs = append(inputs, input)
		}
	}

	output := domain.TxOutput{
		Value:  trans.Amount,
		PubKey: trans.To,
	}
	outputs = append(outputs, output)

	if acc > trans.Amount {
		outpuT := domain.TxOutput{
			Value:  acc - trans.Amount,
			PubKey: trans.From,
		}

		outputs = append(outputs, outpuT)
	}

	tx = domain.Transaction{
		ID:      "",
		Inputs:  inputs,
		Outputs: outputs}

	tx.ID = tx.Hash()

	return tx, nil
}

func (bc blockchain) FindSpendableOutputs(ctx context.Context, trans domain.Trans) (int, map[string][]int, error) {
	unspentOuts := make(map[string][]int)
	accumulated := 0

	unspentTxs, err := bc.FindUnspentTransactions(ctx, trans.From)
	if err != nil {
		return accumulated, unspentOuts, err
	}

	// js, err := json.MarshalIndent(unspentTxs, " ", " ")
	// if err != nil {
	// 	log.Fatalf("failed tx:" + err.Error())
	// }
	// log.Printf("unspentTxs js:\n%s", string(js))

Work:
	for _, tx := range unspentTxs {
		for outIdx, out := range tx.Outputs {
			if bc.canBeUnlocked(ctx, out, trans.From) && accumulated < trans.Amount {
				accumulated += out.Value
				unspentOuts[tx.ID] = append(unspentOuts[tx.ID], outIdx)
				if accumulated >= trans.Amount {
					break Work
				}
			}
		}
	}

	// js, err = json.MarshalIndent(unspentOuts, " ", " ")
	// if err != nil {
	// 	log.Fatalf("failed tx:" + err.Error())
	// }
	// log.Printf("unspentOuts:\n%s", string(js))

	return accumulated, unspentOuts, nil
}
