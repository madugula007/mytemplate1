package repository

import (
	"context"
	"errors"
	
	"reflect"
	//"strings"

	"time"

	"gotemplate/core/domain"
	"gotemplate/core/port"
	"gotemplate/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepository struct {
	db  *DB
	log *logger.Logger
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *DB, log *logger.Logger) *UserRepository {
	return &UserRepository{
		db,
		log,
	}
}

// id, name, email, password, role, created_at, updated_at
type Uservalues struct {
	Id         uint64
	Name       string
	Email      string
	Password   string
	Role       string
	Created_at time.Time
	Updated_at time.Time
}

func generateColumnsFromStruct(instance interface{}) []string {
	var columns []string

	val := reflect.Indirect(reflect.ValueOf(instance))
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("select")
		if tag != "" {
			columns = append(columns, tag)
		}
	}

	return columns
}
func generateMapFromStruct(instance interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.Indirect(reflect.ValueOf(instance))
	typ := val.Type()
	//fmt.Println("struct:", typ)
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("insert")
		if tag != "" {
			result[tag] = val.Field(i).Interface()
		}
	}

	return result
}

// Function for queries returning a single instance
func (ur *UserRepository) executeSingleRowQuery(builder sq.Sqlizer, s interface{}, ctx context.Context) (interface{}, error) {
	val := reflect.Indirect(reflect.ValueOf(s))
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	row := ur.db.QueryRow(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	// Create a new instance of the target struct
	targetInstance := reflect.New(val.Type()).Interface()

	// Create a slice of interfaces to pass to row.Scan
	scanArgs := make([]interface{}, val.NumField())
	for i := range scanArgs {
		scanArgs[i] = reflect.New(val.Type().Field(i).Type).Interface()
	}

	// Scan into the slice of interfaces
	if err := row.Scan(scanArgs...); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil for no rows, without an error
		}
		return nil, err
	}

	// Copy values from the scanArgs to the target struct
	for i := range scanArgs {
		reflect.ValueOf(targetInstance).Elem().Field(i).Set(reflect.ValueOf(scanArgs[i]).Elem())
	}

	return targetInstance, nil
}

func (ur *UserRepository) executeQueries(builder sq.Sqlizer, ctx context.Context) (interface{}, error) {
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		ur.log.Debug("error while running query", err.Error())
		return nil, err
	}
	defer rows.Close()
	//rows.FieldDescriptions()
	columns := rows.FieldDescriptions()
	// if err != nil {
	//     return nil, err
	// }

	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))

		for i := range columns {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})

		for i, col := range columns {
			// Convert the interface{} value to the actual type using type assertions
			value := *values[i].(*interface{})
			rowMap[col.Name] = value
		}

		result = append(result, rowMap)
	}

	return result, nil
}


// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// ur.log.Info("Inside create User")
	// ur.log.Debug("repo name:", user.Name)
	// ur.log.Debug("repo email:", user.Email)
	// ur.log.Debug("repo password:", user.Password)

	query := psql.Insert("users").SetMap(generateMapFromStruct(user)).
	Suffix("returning *")
	u, err := ur.executeSingleRowQuery(query, user, ctx)
	if err != nil {
		return nil, err
	}
	resultSlice, ok := u.(*domain.User)
	if !ok {
		// Handle the case where the type assertion fails
		ur.log.Debug("Type assertion failed")
		
		return nil,errors.New("type assertion failed")
	}
	return resultSlice, nil
}

// GetUserByID gets a user by ID from the database
func (ur *UserRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {

	ur.log.Info("Came inside getuser by id")
	//var user domain.User

	//  query := psql.Select("*").
	// 	From("users").
	// 	Where(sq.Eq{"id": id}).
	// 	Limit(1)

	var u1 domain.User

	columns := generateColumnsFromStruct(u1)

	query := psql.Select(columns...).
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)

	if err != nil {
		ur.log.Info(err.Error())
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])

	// err = ur.db.QueryRow(ctx, sql, args...).Scan(
	// 	&user.ID,
	// 	&user.Name,
	// 	&user.Email,
	// 	&user.Password,
	// 	&user.Role,
	// 	&user.CreatedAt,
	// 	&user.UpdatedAt,
	// )
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, port.ErrDataNotFound
		}
		return nil, err
	}

	return &u, nil
}

// GetUserByEmailAndPassword gets a user by email from the database
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := psql.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers lists all users from the database
func (ur *UserRepository) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	//var user domain.User
	//var users []domain.User

	query := psql.Select("*").
		From("users").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

		

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])

	// 	for rows.Next() {
	// 	err := rows.Scan(
	// 		&user.ID,
	// 		&user.Name,
	// 		&user.Email,
	// 		&user.Password,
	// 		&user.Role,
	// 		&user.CreatedAt,
	// 		&user.UpdatedAt,
	// 	)
	if err != nil {
		return nil, err
	}

	// 	users = append(users, user)
	// }

	return user, nil
}

// UpdateUser updates a user by ID in the database
func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// name := nullString(user.Name)
	// email := nullString(user.Email)
	// password := nullString(user.Password)
	// role := nullString(string(user.Role))

	query := psql.Update("users").
		// Set("name", sq.Expr("COALESCE(?, name)", name)).
		// Set("email", sq.Expr("COALESCE(?, email)", email)).
		// Set("password", sq.Expr("COALESCE(?, password)", password)).
		// Set("role", sq.Expr("COALESCE(?, role)", role)).
		// Set("updated_at", time.Now()).
		SetMap(generateMapFromStruct(user)).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()

	ur.log.Debug("SQL:", sql)
	ur.log.Debug("Args:", args)
	if err != nil {
		ur.log.Debug("Error:", err.Error())
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := psql.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
