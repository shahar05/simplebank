package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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
	log.Printf("config: %v\n", config)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot conne ct to db:", err)
	}
	log.Printf("after sql \n")
	runDBMigration(config.MigrationURL, config.DBSource)
	log.Printf("after runDBMigration \n")
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

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	log.Printf("runDBMigration 1\n")

	if err != nil {
		log.Printf("runDBMigration 2: %v\n", err)
		log.Fatal("cannot create new migration instance: ", err)
	}

	if err = migration.Up(); err != nil {
		log.Printf("runDBMigration 3: %v\n", err)
		log.Fatal("failed to run migrateup:", err)
	}

	log.Println("db migrate successfully")
}
