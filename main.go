package main

import (
	"database/sql"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gitaepark/pha/loader"
	"github.com/gitaepark/pha/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := loader.ConnectDB(config.DBDriver, config.DBSource, config.MigrationURL)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	runServer(config, conn)
}

func runServer(config util.Config, conn *sql.DB) {
	server, err := loader.NewServer(config, conn)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
