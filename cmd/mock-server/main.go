package main

import (
	"flag"
	"fmt"
	"log"

	"protocol-adapter/http-mock-server/config"
	"protocol-adapter/http-mock-server/core"

	"github.com/valyala/fasthttp"
)

func main() {
	appConfig := config.NewAppConf()

	flag.Parse()
	fixturesPaths := flag.Args()

	fixtures, err := core.CollectFixtures(fixturesPaths)
	if err != nil {
		log.Fatalln(err)
	}

	mock := core.NewMock(fixtures)
	log.Printf("Load %v fixtures", len(fixtures))

	handler := core.NewHandler(mock)

	addr := fmt.Sprintf(":%d", appConfig.MockServer.Port)
	log.Printf("Start server on port %v\n", appConfig.MockServer.Port)
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}
