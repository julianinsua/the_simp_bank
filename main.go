package main

import (
	"database/sql"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/julianinsua/the_simp_bank.git/api"
	"github.com/julianinsua/the_simp_bank.git/internal/database"
	"github.com/julianinsua/the_simp_bank.git/util"
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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("failed to initialize server", err)
	}
}
