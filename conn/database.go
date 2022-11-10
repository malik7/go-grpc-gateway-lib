package conn

import (
	"fmt"
	"time"

	// Init the mysql database module
	_ "github.com/go-sql-driver/mysql"

	// Init the postgres database module
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// DBConfig holds per database config
// DSN				: connection string. It usually contains the username, password, and address
// MaxConnOpen		: maximum number of open connections to the database. Then there is no limit on the number of open connections. Default: 0
// MaxConnDuration	: maximum amount of time a connection may be reused. Expired connections may be closed lazily before reuse. If duration <= 0, connections are reused forever. Default value: 0
// MaxConnIdle		: maximum number of connections in the idle connection pool. If MaxConnIdle > 0 && < MaxIdleConns then the MaxConnIdle will be reduced to MaxConnOpen limit. Default value: 0
type DBConfig struct {
	Dialect         string        `envconfig:"DIALECT"`
	Host            string        `envconfig:"HOST"`
	Port            int           `envconfig:"PORT"`
	Name            string        `envconfig:"NAME"`
	Username        string        `envconfig:"USER_NAME"`
	Password        string        `envconfig:"PASSWORD"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE"`
}

// sqlxConnect is used to mock sqlx.Open
var sqlxConnect = sqlx.Connect

// InitDB init the database from config to database connection
func InitDB(cfg DBConfig) (*sqlx.DB, error) {
	dsn, err := cfg.toDSN()
	if err != nil {
		return nil, err
	}

	db, err := sqlxConnect(cfg.Dialect, dsn)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxConnOpen)
	db.SetMaxIdleConns(cfg.MaxConnIdle)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)

	return db, nil
}

func (cfg DBConfig) toDSN() (string, error) {
	switch cfg.Dialect {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name), nil
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name), nil
	default:
		return "", fmt.Errorf("Dialect is not supported. expected: (msql|postgres), got: %s", cfg.Dialect)
	}
}
