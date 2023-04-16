package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/shahar05/simplebank/api"
	db "github.com/shahar05/simplebank/db/sqlc"
	"github.com/shahar05/simplebank/util"
)

// TODO: try to remove 	_ "github.com/lib/pq" and see what happen. Hint:(Cant talk to to PQ server...)
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}