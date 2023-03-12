package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/thehaung/simplebank/api"
	"github.com/thehaung/simplebank/config"
	db "github.com/thehaung/simplebank/db/sqlc"
	"log"
)

func main() {
	conf, err := config.Parse(".")
	if err != nil {
		log.Fatal("main - config.Parse. Error:", err)
	}
	conn, err := sql.Open(conf.DbDriver, conf.DbAddress)
	if err != nil {
		log.Fatal("main - sql.Open. Error:", err)
	}

	dbStore := db.NewStore(conn)
	httpServer, err := api.NewHttpServer(conf, dbStore)
	if err != nil {
		log.Fatal("main - api.NewHttpServer(). Error:", err)
	}

	err = httpServer.Start(conf.HttpServerAddress)
	if err != nil {
		log.Fatal("main - httpServer.Start(). Error:", err)
	}
}
