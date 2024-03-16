package repository

import (
	"context"
	"errors"
	"fmt"

	//"os"
	"gotemplate/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/**
 * DB is a wrapper for PostgreSQL database connection
 * that uses pgxpool as database driver
 */
type DB struct {
	*pgxpool.Pool
}

type DBInterface interface {
	Close()
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error, levels ...pgx.TxIsoLevel) error
	ReadTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	// You might have other methods that you want to expose through the interface
}

var _ DBInterface = (*DB)(nil)

// NewDB creates a new PostgreSQL database instance
func NewDB(ctx context.Context, c config.Econfig) (*DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s  sslmode=disable",
		c.DBUsername(),
		c.DBPassword(),
		c.DBHost(),
		c.DBPort(),
		c.DBDatabase(),
	
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	
	config.MaxConns = int32(c.MaxConns())                                     // Maximum number of connections in the pool.
	config.MinConns = int32(c.MinConns())                                     // Minimum number of connections to keep in the pool.
	config.MaxConnLifetime = time.Duration(c.MaxConnLifetime()) * time.Minute // Maximum lifetime of a connection.
	config.MaxConnIdleTime = time.Duration(c.MaxConnIdleTime()) * time.Minute // Maximum idle time of a connection in the pool.

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

func (db *DB) WithTx(ctx context.Context, fn func(tx pgx.Tx) error, levels ...pgx.TxIsoLevel) error {
	var level pgx.TxIsoLevel
	if len(levels) > 0 {
		level = levels[0]
	} else {
		level = pgx.ReadCommitted // Default value
	}
	return db.inTx(ctx, level, "", fn)
}

func (db *DB) ReadTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return db.inTx(ctx, pgx.ReadCommitted, pgx.ReadOnly, fn)

}

func (db *DB) inTx(ctx context.Context, level pgx.TxIsoLevel, access pgx.TxAccessMode,
	fn func(tx pgx.Tx) error) (err error) {

	conn, errAcq := db.Pool.Acquire(ctx)
	if errAcq != nil {
		return fmt.Errorf("acquiring connection: %w", errAcq)
	}
	defer conn.Release()

	opts := pgx.TxOptions{
		IsoLevel:   level,
		AccessMode: access,
	}

	tx, errBegin := conn.BeginTx(ctx, opts)
	if errBegin != nil {
		return fmt.Errorf("begin tx: %w", errBegin)
	}

	defer func() {
		errRollback := tx.Rollback(ctx)
		if !(errRollback == nil || errors.Is(errRollback, pgx.ErrTxClosed)) {
			err = errRollback
		}
	}()

	if err := fn(tx); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			return fmt.Errorf("rollback tx: %v (original: %w)", errRollback, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
