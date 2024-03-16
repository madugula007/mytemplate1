package repository

import (
	"context"
	"errors"
	"time"

	//"encoding/json"

	//"reflect"
	//"strings"

	"gotemplate/core/domain"
	//"gotemplate/core/port"
	"gotemplate/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepository struct {
	Db  *DB
	log *logger.Logger
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(Db *DB, log *logger.Logger) *UserRepository {
	return &UserRepository{
		Db,
		log,
	}
}

// CreateUser creates a new user in the database
// func (ur *UserRepository) CreateUser(gctx *gin.Context, user *domain.User) (*domain.UserDB, error) {
func (ur *UserRepository) CreateUser(gctx *gin.Context, user *domain.User) (domain.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//ur.log.Debug("USer:", user)

	query := psql.Insert("users").SetMap(generateMapFromStruct(user, "insert")).Suffix("returning *")
	// sql, args, err := query.ToSql()
	// ur.log.Debug("sql:", sql)
	// ur.log.Debug("args:", args)
	// if err != nil {
	// 	ur.log.Debug("Error:", err.Error())
	// 	return err
	// }

	// e, er := ur.Db.Exec(ctx, sql, args...)
	// ur.log.Debug("e:", e)

	p, err := InsertReturning(ctx, ur.Db, query, pgx.RowToStructByName[domain.User], ur.log)
	// ur.log.Debug("P:", p.Insert())
	// ur.log.Debug("error:", err)
	//ur.log.Debug("user:", p)
	return p, err
	//return Insert(ctx, ur.Db, query, pgx.RowToAddrOfStructByPos[domain.UserDB], ur.log)
}

// GetUserByID gets a user by ID from the database
func (ur *UserRepository) GetUserByID(gctx *gin.Context, id uint64) (*domain.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ur.log.Info("Came inside getuser by id")
	var u1 domain.UserDB
	columns := generateColumnsFromStruct(u1, "select")
	query := psql.Select(columns...).
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	return SelectOneOK(ctx, ur.Db, query, pgx.RowToAddrOfStructByNameLax[domain.User], ur.log)
}

// GetUserByEmailAndPassword gets a user by email from the database
func (ur *UserRepository) GetUserByEmail(gctx *gin.Context, email string) (*domain.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//var user domain.User
	query := psql.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)
	return SelectOneOK(ctx, ur.Db, query, pgx.RowToAddrOfStructByName[domain.User], ur.log)
}

// ListUsers lists all users from the database
func (ur *UserRepository) ListUsers(gctx *gin.Context, skip, limit uint64) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := psql.Select("name,email,password").
		From("users").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)
	//return SelectRows(ctx, ur.Db, query, pgx.RowToAddrOfStructByNameLax[domain.User], ur.log)

	// sql, args, err := query.ToSql()
	// if err != nil {
	// 	ur.log.Debug("Logs :", err.Error())
	// 	return nil, err
	// }
	//s,_:=SelectWithSecondary(ctx,ur.Db,"select",domain.User{},sql,args)

	//return nil,nil
	//return SelectRows(ctx, ur.Db, query, RowToStructByTagExp[domain.User], ur.log)
	if ctx.Err() == context.DeadlineExceeded {
		ur.log.Error("Context deadline exceeded")
		return nil, errors.New("context deadline exceeded")
	}
	return SelectRowsTag[domain.User](ctx, ur.Db, query, ur.log, "select")

}

// UpdateUser updates a user by ID in the database
func (ur *UserRepository) UpdateUser(gctx *gin.Context, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := psql.Update("users").
		// Set("name", sq.Expr("COALESCE(?, name)", name)).
		// Set("email", sq.Expr("COALESCE(?, email)", email)).
		// Set("password", sq.Expr("COALESCE(?, password)", password)).
		// Set("role", sq.Expr("COALESCE(?, role)", role)).
		// Set("updated_at", time.Now()).
		SetMap(generateMapFromStruct(user, "insert")).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")
	return UpdateReturning(ctx, ur.Db, query, pgx.RowToAddrOfStructByPos[domain.User], ur.log)

}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(gctx *gin.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := psql.Delete("users").
		Where(sq.Eq{"id": id})
	_, err := Delete(ctx, ur.Db, query, ur.log)
	return err
}
