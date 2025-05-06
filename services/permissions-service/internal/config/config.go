package config

import (
	"github.com/alecthomas/kong"
)

type Config struct {
	KeyCloakURL    string `help:"KeyCloak base url" env:"KEYCLOAK_URL" required:"true"`
	KeyCloakRealm  string `help:"KeyCloak Realm name" env:"KEYCLOAK_REALM" required:"true"`
	KeyCloakClient string `help:"KeyCloak client id" env:"KEYCLOAK_CLIENT" required:"true"`
	KeyCloakSecret string `help:"KeyCloak client secret" env:"KEYCLOAK_SECRET" required:"true"`
	DatabaseURL    string `help:"DB connection string" env:"DATABASE_URL" required:"true"`
	Port           string `help:"Port that grpc will listen on" env:"PORT" default:"10000"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	parser, err := kong.New(cfg)
	if err != nil {
		return nil, err
	}

	// Parse command-line flags, environment variables, and config file
	_, err = parser.Parse(nil)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
