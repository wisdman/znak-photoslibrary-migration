package main

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
		MaxConnections: 30,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetAuthor(title string) (string, error) {
	var author string
	if err := DB.QueryRow(`SELECT author FROM album WHERE name = $1`, title).Scan(&author); err != nil {
		return "", err
	}

	return author, nil
}