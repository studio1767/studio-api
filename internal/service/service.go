package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"

	"github.com/parlaynu/studio1767-api/api/graph/model"
	"github.com/parlaynu/studio1767-api/internal/config"
)

type Service interface {
	CreateProject(ctx context.Context, p *model.NewProject) (*model.Project, error)
	Projects(ctx context.Context) ([]*model.Project, error)
	ProjectById(ctx context.Context, id string) (*model.Project, error)
}

func New(cfg *config.Config) (Service, error) {
	// create the db config
	dbCfg := mysql.NewConfig()
	dbCfg.Net = "tcp"
	dbCfg.Addr = fmt.Sprintf("%s:%d", cfg.Db.Server, cfg.Db.Port)
	dbCfg.DBName = cfg.Db.DbName
	dbCfg.User = cfg.Db.UserName
	dbCfg.Passwd = cfg.Db.Password

	if len(cfg.Db.CaCertFile) > 0 {
		pem, err := os.ReadFile(cfg.Db.CaCertFile)
		if err != nil {
			return nil, err
		}
		certs := x509.NewCertPool()
		certs.AppendCertsFromPEM(pem)
		tlsCfg := &tls.Config{
			RootCAs: certs,
		}
		mysql.RegisterTLSConfig("maria", tlsCfg)

		dbCfg.TLSConfig = "maria"
	}

	fmt.Println(dbCfg.FormatDSN())

	// connect to the database
	var err error
	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to read ca certificate file: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	svc := service{
		db: db,
	}

	return &svc, nil
}

type service struct {
	db *sql.DB
}
