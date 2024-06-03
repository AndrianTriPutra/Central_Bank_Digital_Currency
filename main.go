package main

import (
	"atp/cbdc/app/endpoint"
	"atp/cbdc/app/usecase/utransaction"
	"atp/cbdc/pkg/adapter/levelDB"
	"atp/cbdc/pkg/repository/blockchain"
	"atp/cbdc/pkg/repository/crud"
	"atp/cbdc/pkg/repository/pow"

	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	log.Println(" ========================== CBDC ========================== ")
	log.Println("created by : atp")
	log.Println("version : v1.0 2024-06-03")
	log.Println(" ========================== CBDC ========================== ")
}

func main() {
	ctx := context.Background()

	setting := pow.Setting{
		Difficult: 13,
	}
	genesis := blockchain.Genesis{
		Address: "Bank Central",
		Balance: 100,
	}

	db, err := levelDB.NewConnection("./database")
	if err != nil {
		log.Fatalf("failed connect to database:" + err.Error())
	}

	//init repository
	repoPoW := pow.NewRepository(setting)
	repoCrud := crud.NewRepository(db)
	repoBC := blockchain.NewRepository(genesis, repoCrud, repoPoW)

	//initialization blockchain
	err = repoBC.InitBlockchain(ctx)
	if err != nil {
		log.Fatalf("failed InitBlockchain:" + err.Error())
	}

	//init usecase
	ucase := utransaction.NewUsecase(repoBC, repoCrud)

	echoNew := echo.New()
	echoNew.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logEcho := fmt.Sprintf("{status:%v} {method:%v} {latency:%v} {uri:%v} {user_agent:%v}", values.Status, values.Method, values.Latency, values.URI, values.UserAgent)
			if values.Status != 200 {
				log.Printf("[error] [logEcho] %s", logEcho)
			} else {
				log.Printf("[info] [logEcho] %s", logEcho)
			}
			return nil
		},
	}))

	//endpoint
	endpoint.NewHandler(echoNew, "atp/cbdc/", ucase)

	errServer := make(chan error)
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	echoNew.Server.TLSConfig = cfg
	echoNew.Server.Addr = ":8008"
	//optional
	timeout := 10 * time.Minute
	echoNew.Server.ReadTimeout = timeout
	echoNew.Server.WriteTimeout = timeout
	echoNew.Server.IdleTimeout = timeout

	runServer := func() {
		log.Printf("[info] server running on port [%s]", echoNew.Server.Addr)
		errServer <- echoNew.Server.ListenAndServe()
	}

	go runServer()

	for {
		select {
		case <-ctx.Done():
			ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			go func(ctx context.Context) {
				defer cancel()
				// shutdown server
				if err := echoNew.Shutdown(ctxShutDown); err != nil {
					log.Fatalf("[fatal] server shutdown failed:%s" + err.Error())
				}
				log.Fatal("[fatal] server exited properlys")
			}(ctx)

		case err := <-errServer:
			log.Fatalf("[fatal] server error got:%s" + err.Error())
		}
	}

}
