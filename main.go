package main

import (
	"MelBank/api"
	db "MelBank/db/sqlc"
	"MelBank/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// @title Bank App API
// @version 1.0
// @description Api for Bank app

// @host localhost:8080
// @BasePath /

// @securityDefinition.apikey ApiKeyAuth
// @in header
// @name Bearer Auth

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db:", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
