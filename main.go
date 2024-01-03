package main

import (
	"database/sql"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/julianinsua/the_simp_bank/api"
	"github.com/julianinsua/the_simp_bank/internal/database"
	"github.com/julianinsua/the_simp_bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to log config file: ", err)
	}

	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("unable to create database connection", err)
	}

	store := database.NewStore(db)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create new server", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("failed to initialize server", err)
	}
}
