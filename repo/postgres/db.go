package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

/**
 * DB is a wrapper for PostgreSQL database connection
 * that uses pgxpool as database driver
 */
type DB struct {
	*pgxpool.Pool
}

// NewDB creates a new PostgreSQL database instance
func NewDB(ctx context.Context) (*DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable pool_max_conns=20",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	/* Remember to change these values according to our need and placing these
			values in config.yaml
			Longer MaxConnIdleTime values can result in connections being kept open for a longer duration, potentially utilizing more resources
			If application has consistent and steady traffic .. consider for shorter MaxConnIdleTime
		MaxConnLifetime:maximum amount of time a connection can be open before the pool considers it eligible for closing
	Longer MaxConnLifetime values allow connections to be reused for a more extended period, reducing the overhead of establishing new connections.
	Extremely long MaxConnLifetime is not preferable.Longer MaxConnLifetime utilises more resources.
	Adjust all parameters in production after monitoring application's performance.

	*/
	config.MaxConns = 10                      // Maximum number of connections in the pool.
	config.MinConns = 3                       // Minimum number of connections to keep in the pool.
	config.MaxConnLifetime = 10 * time.Minute // Maximum lifetime of a connection.
	config.MaxConnIdleTime = 3 * time.Minute  // Maximum idle time of a connection in the pool.

	db, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	db.Pool.Close()
}
