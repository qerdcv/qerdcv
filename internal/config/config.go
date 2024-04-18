package config

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const AppName = "qerdcv"

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Addr            string        `default:":8080"`
	ReadTimeout     time.Duration `split_words:"true" default:"3s"`
	WriteTimeout    time.Duration `split_words:"true" default:"5s"`
	ShutdownTimeout time.Duration `split_words:"true" default:"15s"`
}

type DBConfig struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
}

func (c DBConfig) DSN() string {
	u := &url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.Name,
		User:     url.UserPassword(c.Username, c.Password),
		RawQuery: "sslmode=disable",
	}

	return u.String()
}

func New() (Config, error) {
	var c Config
	if err := envconfig.Process(AppName, &c); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return c, nil
}
