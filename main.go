package main

import (
	"database/sql"
	"log"

	"github.com/Xinnan-Alex/simplebank/api"
	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/Xinnan-Alex/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start sercer:", err)
	}
}
