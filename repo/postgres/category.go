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
 * CategoryRepository implements port.CategoryRepository interface
 * and provides an access to the postgres database
 */
type CategoryRepository struct {
	Db *DB
	log *logger.Logger
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(Db *DB,log *logger.Logger) *CategoryRepository {
	return &CategoryRepository{
		Db,
		log,
	}
}

// CreateCategory creates a new category record in the database
func (cr *CategoryRepository) CreateCategory(gctx *gin.Context, category *domain.Category) (*domain.Category, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := psql.Insert("categories").
		Columns("name").
		Values(category.Name).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

cr.Db.QueryRow(ctx, sql, args...)


	err = cr.Db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category record from the database by id
func (cr *CategoryRepository) GetCategoryByID(gctx *gin.Context, id uint64) (*domain.Category, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var category domain.Category

	query := psql.Select("*").
		From("categories").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.Db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, port.ErrDataNotFound
		}
		return nil, err
	}

	return &category, nil
}

// ListCategories retrieves a list of categories from the database
func (cr *CategoryRepository) ListCategories(gctx *gin.Context, skip, limit uint64) ([]domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var category domain.Category
	var categories []domain.Category

	query := psql.Select("*").
		From("categories").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategory updates a category record in the database
func (cr *CategoryRepository) UpdateCategory(gctx *gin.Context, category *domain.Category) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := psql.Update("categories").
		Set("name", category.Name).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": category.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.Db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category record from the database by id
func (cr *CategoryRepository) DeleteCategory(gctx *gin.Context, id uint64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := psql.Delete("categories").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = cr.Db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
