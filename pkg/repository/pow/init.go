package pow

import (
	"atp/cbdc/pkg/utils/domain"
	"math/big"
)

type ProofOfWork struct {
	Block  *domain.Block
	Target *big.Int
}

type Setting struct {
	Difficult int
}

type repository struct {
	setting Setting
}

func NewRepository(setting Setting) RepositoryI {
	return repository{
		setting,
	}
}

type RepositoryI interface {
	NewProof(b *domain.Block) *ProofOfWork
	Run(pOw *ProofOfWork) (uint, []byte)
	Validate(pOw *ProofOfWork, nonce uint) bool
}
