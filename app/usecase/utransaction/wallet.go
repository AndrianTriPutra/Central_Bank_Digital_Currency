package utransaction

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
	"encoding/json"
)

func (u usecase) GetBalance(ctx context.Context, addr string) (int, error) {
	balance := 0

	utxo, err := u.repoBC.FindUTXO(ctx, addr)
	if err != nil {
		return 0, err
	}

	for _, out := range utxo {
		balance += out.Value
	}

	return balance, nil
}

func (u usecase) GetWallet(ctx context.Context) ([]domain.Wallet, error) {
	var wallet []domain.Wallet
	w, err := u.repoCrud.GetWallet(ctx)
	if err != nil {
		return wallet, err
	}

	for i, val := range w {
		x := domain.Wallet{
			Index:   i + 1,
			Address: val,
		}
		wallet = append(wallet, x)
	}

	return wallet, nil
}

func (u usecase) CheckWallet(ctx context.Context, addr string) ([]string, bool, error) {
	w, err := u.repoCrud.GetWallet(ctx)
	if err != nil {
		return nil, false, err
	}

	found := false
	for _, val := range w {
		if addr == val {
			found = true
		}
	}

	return w, found, nil
}

func (u usecase) UpdateWallet(ctx context.Context, addr []string) error {
	w := domain.DWallet{
		Address: addr,
	}

	js, err := json.Marshal(w)
	if err != nil {
		return err
	}

	err = u.repoCrud.PUT(ctx, "wallet", string(js))
	if err != nil {
		return err
	}

	return nil
}

func (u usecase) History(ctx context.Context, addr string) ([]domain.Block, error) {
	var resp []domain.Block
	data, err := u.repoCrud.GetChain(ctx)
	if err != nil {
		return nil, err
	}

	for _, val := range data {
		_, err = json.Marshal(val)
		if err != nil {
			return nil, err
		}

		save := false
		for _, trans := range val.Transaction {
			for _, input := range trans.Inputs {
				if input.Sig == addr {
					save = true
					break
				}
			}

			for _, output := range trans.Outputs {
				if output.PubKey == addr {
					save = true
					break
				}
			}
		}

		if save {
			resp = append(resp, val)
		}
	}
	return resp, nil
}
