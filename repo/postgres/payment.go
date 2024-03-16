package repository

import (
	"context"
	"time"

	"gotemplate/core/domain"
	"gotemplate/core/port"
	"gotemplate/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

/**
 * PaymentRepository implements port.PaymentRepository interface
 * and provides an access to the postgres database
 */
type PaymentRepository struct {
	Db  *DB
	log *logger.Logger
}

// NewPaymentRepository creates a new payment repository instance
func NewPaymentRepository(Db *DB, log *logger.Logger) *PaymentRepository {
	return &PaymentRepository{
		Db,
		log,
	}
}

// CreatePayment creates a new payment record in the database
func (pr *PaymentRepository) CreatePayment(gctx *gin.Context, payment *domain.Payment) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := psql.Insert("payments").
		Columns("name", "type", "logo").
		Values(payment.Name, payment.Type, payment.Logo).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.Db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.Name,
		&payment.Type,
		&payment.Logo,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentByID retrieves a payment record from the database by id
func (pr *PaymentRepository) GetPaymentByID(gctx *gin.Context, id uint64) (*domain.Payment, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var payment domain.Payment

	query := psql.Select("*").
		From("payments").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.Db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.Name,
		&payment.Type,
		&payment.Logo,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, port.ErrDataNotFound
		}
		return nil, err
	}

	return &payment, nil
}

// ListPayments retrieves a list of payments from the database
func (pr *PaymentRepository) ListPayments(gctx *gin.Context, skip, limit uint64) ([]domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payment domain.Payment
	var payments []domain.Payment

	query := psql.Select("*").
		From("payments").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pr.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&payment.ID,
			&payment.Name,
			&payment.Type,
			&payment.Logo,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

// UpdatePayment updates a payment record in the database
func (pr *PaymentRepository) UpdatePayment(gctx *gin.Context, payment *domain.Payment) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	name := nullString(payment.Name)
	paymentType := nullString(string(payment.Type))
	logo := nullString(payment.Logo)

	query := psql.Update("payments").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("type", sq.Expr("COALESCE(?, type)", paymentType)).
		Set("logo", sq.Expr("COALESCE(?, logo)", logo)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": payment.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.Db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.Name,
		&payment.Type,
		&payment.Logo,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// DeletePayment deletes a payment record from the database by id
func (pr *PaymentRepository) DeletePayment(gctx *gin.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := psql.Delete("payments").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pr.Db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
