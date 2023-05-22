package db

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/studio1767/studio-api/internal/config"
)

func NewClient(cfg *config.Config, tlsConfig *tls.Config) (*sql.DB, error) {
	// create the db config
	dbConfig := mysql.NewConfig()
	dbConfig.Net = "tcp"
	dbConfig.Addr = fmt.Sprintf("%s:%d", cfg.Db.Server, cfg.Db.Port)
	dbConfig.DBName = cfg.Db.DbName
	dbConfig.User = cfg.Db.UserName
	dbConfig.Passwd = cfg.Db.Password

	err := mysql.RegisterTLSConfig("maria", tlsConfig)
	if err != nil {
		return nil, err
	}
	dbConfig.TLSConfig = "maria"

	// connect to the database
	client, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return client, nil
}
