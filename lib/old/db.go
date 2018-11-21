package old

import (
	"log"

	"github.com/jackc/pgx"
)

var DB *pgx.ConnPool

func init() {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		log.Fatal(err)
	}

	DB, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Progress() {
	var migrated uint64
	var all uint64

	err := DB.QueryRow(`
    SELECT (SELECT COUNT(*) FROM album WHERE migrated) as migrated, (SELECT COUNT(*) FROM album) as all
  `).Scan(&migrated, &all)

	if err != nil {
		log.Fatal(err)
	}

	percentage := (migrated / all) * 100

	log.Printf("=== Migrated %d%% | %d of %d ===", percentage, migrated, all)
}
