package utransaction_test

import (
	"atp/cbdc/app/usecase/utransaction"
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

func Test_Balance(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	//log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "app")]
	path := base + "database/"
	//log.Printf("path:%s", path)

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
	ucase := utransaction.NewUsecase(repoBC, repoCrud)

	addr := "Bank Central"
	//addr := "Bank D"
	balance, err := ucase.GetBalance(ctx, addr)
	if err != nil {
		log.Fatalf("failed GetBalance:" + err.Error())
	}
	log.Printf("balance from addr [%s] is [%v]", addr, balance)
}
