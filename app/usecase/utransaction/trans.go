package utransaction

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
)

func (u usecase) Transaction(ctx context.Context, trans domain.Trans) error {
	err := u.repoBC.Send(ctx, trans)
	if err != nil {
		return err
	}

	return nil
}
