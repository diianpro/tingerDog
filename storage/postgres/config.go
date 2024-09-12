package postgres

import "fmt"

type Config struct {
	DataBase     string `config:"POSTGRES_DB,required"`
	User         string `config:"POSTGRES_USER,required"`
	Password     string `config:"POSTGRES_PASSWORD,required"`
	Host         string `config:"POSTGRES_HOST,required"`
	Port         string `config:"POSTGRES_PORT,required"`
	MaxOpenConns int32  `config:"POSTGRES_MAX_OPEN_CONNS,required"`
	MaxIdleConns int32  `config:"POSTGRES_MAX_IDLE_CONNS,required"`
}

// ConnectionString returns a connection string for the database.
// if databaseName is empty, database name will be taken from the environment variable POSTGRES_DATABASE.
func ConnectionString(cfg *Config) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", //nolint:nosprintfhostport
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DataBase,
	)
}
