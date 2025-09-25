package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Config represents the database configuration
type Config struct {
	Name        string `env:"APP_NAME,default=go-learning-db"`
	Environment string
	Database    string `env:"DB_NAME,required"`
	DBHost      string `env:"DB_HOST,required"`
	DBPort      string `env:"DB_PORT,required"`
	DBUser      string `env:"DB_USER,required"`
	DBSecret    string `env:"DB_PASS,required"`
	SSLMode     string `env:"DB_SSL_MODE,default=disable"`
	Metrics     *prometheus.Registry
}

// Connection represents a database connection that implements the db interface
type Connection struct {
	name        string
	environment string
	pool        *pgxpool.Pool
}

// NewConnection creates a new database connection
func NewConnection(cfg Config) (*Connection, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBSecret, cfg.DBHost, cfg.DBPort, cfg.Database, cfg.SSLMode)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &Connection{
		name:        cfg.Name,
		environment: cfg.Environment,
		pool:        pool,
	}, nil
}

// Start initializes the database connection
func (c *Connection) Start() error {
	if c.environment != "test" {
		if err := c.pool.Ping(context.Background()); err != nil {
			zap.L().Error("database connection failed",
				zap.String("name", c.name),
				zap.String("environment", c.environment),
				zap.Error(err))
			return err
		}
		zap.L().Info("database connection established",
			zap.String("name", c.name),
			zap.String("environment", c.environment))
	}
	return nil
}

// Stop closes the database connection
func (c *Connection) Stop() error {
	if c.pool != nil {
		c.pool.Close()
	}
	return nil
}

// Name returns the service name
func (c *Connection) Name() string {
	return c.name
}

// Query executes a query that returns rows
func (c *Connection) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (c *Connection) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query without returning any rows
func (c *Connection) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, arguments...)
}

// BeginTx starts a transaction
func (c *Connection) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, txOptions)
}

// validate validates the configuration
func (cfg *Config) validate() error {
	var errs []error

	if cfg.Name == "" {
		errs = append(errs, fmt.Errorf("db: Name is required"))
	}
	if cfg.Environment == "" {
		errs = append(errs, fmt.Errorf("db: Environment is required"))
	}
	if cfg.Database == "" {
		errs = append(errs, fmt.Errorf("db: Database is required"))
	}
	if cfg.DBHost == "" {
		errs = append(errs, fmt.Errorf("db: DBHost is required"))
	}
	if cfg.DBPort == "" {
		errs = append(errs, fmt.Errorf("db: DBPort is required"))
	}
	if cfg.DBUser == "" {
		errs = append(errs, fmt.Errorf("db: DBUser is required"))
	}
	if cfg.DBSecret == "" {
		errs = append(errs, fmt.Errorf("db: DBSecret is required"))
	}
	if cfg.Metrics == nil {
		errs = append(errs, fmt.Errorf("db: Metrics is required"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %v", errs)
	}
	return nil
}
