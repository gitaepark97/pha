package loader

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ConnectDB(dbDriver, dbSource, migrationURL string) (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	// DB 대기 커넥션 수
	conn.SetMaxIdleConns(2)
	// DB 최대 커넥션 수
	conn.SetMaxOpenConns(2)

	runMigration(migrationURL, dbDriver, dbSource)

	return conn, nil
}

func runMigration(migrationURL, dbDriver, dbSource string) {
	migration, err := migrate.New(migrationURL, dbDriver+"://"+dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migrate instace")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}
