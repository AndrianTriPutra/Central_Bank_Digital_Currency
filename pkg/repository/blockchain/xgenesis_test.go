package blockchain_test

import (
	"atp/cbdc/pkg/adapter/levelDB"
	"atp/cbdc/pkg/repository/blockchain"
	"atp/cbdc/pkg/repository/crud"
	"atp/cbdc/pkg/repository/pow"
	"context"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_GENESIS(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

	ctx := context.Background()

	db, err := levelDB.NewConnection(path)
	if err != nil {
		log.Fatalf("failed connect to database:" + err.Error())
	}

	setting := pow.Setting{
		Difficult: 13,
	}
	genesis := blockchain.Genesis{
		Address: "Bank Central",
		Balance: 100,
	}
	repoPoW := pow.NewRepository(setting)

	repoCrud := crud.NewRepository(db)
	repoBC := blockchain.NewRepository(genesis, repoCrud, repoPoW)
	err = repoBC.InitBlockchain(ctx)
	if err != nil {
		log.Fatalf("failed InitBlockchain:" + err.Error())
	}
}
