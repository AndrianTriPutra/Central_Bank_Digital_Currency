package blockchain_test

import (
	"atp/cbdc/pkg/adapter/levelDB"
	"atp/cbdc/pkg/repository/blockchain"
	"atp/cbdc/pkg/repository/crud"
	"atp/cbdc/pkg/repository/pow"
	"atp/cbdc/pkg/utils/domain"
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_FindSpendableOutputs(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

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

	trans := domain.Trans{
		From:   "Bank Central",
		To:     "Bank A",
		Amount: 10,
	}
	accumulated, unspentOuts, err := repoBC.FindSpendableOutputs(ctx, trans)
	if err != nil {
		log.Fatalf("Test_FindSpendableOutputs:%s", err.Error())
	}
	log.Printf("accumulated:%v", accumulated)
	for key, val := range unspentOuts {
		log.Printf("key [%s] value [%v]", key, val)
	}
}

func Test_NewTrans(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

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

	trans := domain.Trans{
		From:   "Bank Central",
		To:     "Bank A",
		Amount: 10,
	}
	tranS, err := repoBC.NewTransaction(ctx, trans)
	if err != nil {
		log.Fatalf("failed NewTransactione:" + err.Error())
	}

	js, err := json.MarshalIndent(tranS, " ", " ")
	if err != nil {
		log.Fatalf("failed NewTransactione:" + err.Error())
	}
	log.Printf("js:\n%s", string(js))
}

func Test_Send(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

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

	trans := domain.Trans{
		From:   "Bank Central",
		To:     "Bank C",
		Amount: 1,
	}

	err = repoBC.Send(ctx, trans)
	if err != nil {
		log.Fatalf("failed Send:" + err.Error())
	}
}
