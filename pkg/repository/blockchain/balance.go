package blockchain

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
	"encoding/json"
	"strings"
)

func (bc *blockchain) FindUnspentTransactions(ctx context.Context, address string) ([]domain.Transaction, error) {
	var unspentTxs []domain.Transaction
	spentTXOs := make(map[string][]int)

	value, err := bc.repoCrud.GET(ctx, "lh")
	if err != nil {
		return unspentTxs, err
	}
	lastkey := value

	for {
		value, err = bc.repoCrud.GET(ctx, lastkey)
		if err != nil {
			buff := err.Error()
			if strings.Contains(buff, "not found") {
				break
			} else {
				return unspentTxs, err
			}
		}

		if len(value) > 0 {
			var block domain.Block
			err = json.Unmarshal([]byte(value), &block)
			if err != nil {
				return unspentTxs, err
			}

			lastkey = block.Header.PrevHash
			for _, tx := range block.Transaction {
			Outputs:
				for outIdx, out := range tx.Outputs {
					if spentTXOs[tx.ID] != nil {
						for _, spentOut := range spentTXOs[tx.ID] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}

					if bc.canBeUnlocked(ctx, out, address) {
						unspentTxs = append(unspentTxs, *tx)
					}
				}

				if !bc.isGenesis(ctx, tx) {
					for _, in := range tx.Inputs {
						if bc.canUnlock(ctx, in, address) {

							spentTXOs[in.ID] = append(spentTXOs[in.ID], in.Out)
						}
					}
				}
			}

			if block.Header.PrevHash == "0000000000000000000000000000000000000000000000000000000000000000" {
				return unspentTxs, nil
			}
		}
	}

	return unspentTxs, nil
}

func (bc *blockchain) FindUTXO(ctx context.Context, address string) ([]domain.TxOutput, error) {
	var UTXOs []domain.TxOutput
	unspentTrans, err := bc.FindUnspentTransactions(ctx, address)
	if err != nil {
		return UTXOs, err
	}

	for _, tx := range unspentTrans {
		for _, out := range tx.Outputs {
			if bc.canBeUnlocked(ctx, out, address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs, nil
}
