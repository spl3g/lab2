package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spl3g/lab2/internal/proxyproto"
	"github.com/spl3g/lab2/internal/userdb"
	"github.com/spl3g/lab2/services/permissions-service/internal/config"
	"github.com/spl3g/lab2/services/permissions-service/internal/keycloak"
)

type Service struct {
	proxyproto.UnimplementedCentrifugoProxyServer
	conn     *pgxpool.Pool
	storage  *userdb.Queries
	config   *config.Config
	kcClient *keycloak.KCClient
}

func New(config *config.Config) (*Service, error) {
	connCfg, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), connCfg)
	if err != nil {
		return nil, err
	}

	kcClient := keycloak.New(config.KeyCloakURL, config.KeyCloakRealm, config.KeyCloakClient, config.KeyCloakSecret)

	return &Service{
		conn:     conn,
		storage:  userdb.New(conn),
		config:   config,
		kcClient: kcClient,
	}, nil
}
