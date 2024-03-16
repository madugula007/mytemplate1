package domain

import (

	//"time"
	//	"github.com/templatedop/githubrepo/dtime"

	//"database/sql"

	//"github.com/volatiletech/null"
	"github.com/guregu/null"
	"github.com/guregu/null/zero"

	//"github.com/aarondl/opt/null"
	"github.com/jackc/pgtype"
	//"github.com/jackc/pgtype"
	//"github.com/jackc/pgx/v5/pgtype"
)

// UserRole is an enum for user's role
type UserRole string

// UserRole enum values
const (
	Admin   UserRole = "admin"
	Cashier UserRole = "cashier"
)

// id,name,email,password,updated_at,created_time
// User is an entity that represents a user
type User struct {
	ID       uint64 `json:"id" db:"id" select:"-" `
	Name     string `json:"name" insert:"name" select:"name" insert_pickup:"name"`
	Email    string `json:"email" insert:"email" select:"email"`
	Password string `json:"password" insert:"password" select:"password"`
	//Role        UserRole    `json:"role" db:"role" `

	CreatedAt   zero.String `json:"created_at" db:"created_at" insert:"created_at" select:"-" `
	UpdatedAt   null.String `json:"updated_at" db:"updated_at" select:"-" `
	CreatedTime zero.String `json:"created_time" db:"created_time" insert:"created_time" select:"-"`
}

type UserDB struct {
	ID       uint64 `json:"id" db:"id" select:"id" `
	Name     string `json:"name" insert:"name" select:"name" insert_pickup:"name"`
	Email    string `json:"email" insert:"email" select:"email"`
	Password string `json:"password" insert:"password" select:"password"`
	//Role      UserRole  `json:"role" db:"role" select:"role"`
	//CreatedAt   null.Val[string] `json:"created_at" db:"created_at" insert:"created_at"  `
	CreatedAt   null.String `json:"created_at" db:"created_at" select:"created_at" `
	UpdatedAt   null.String `json:"updated_at" db:"updated_at" `
	CreatedTime null.String `db:"created_time"  select:"created_time"`
}

type RegisterRequest struct {
	Name     pgtype.Text `json:"name" validate:"required,min=5" u:"N1" db:"name" example:"John Doe"`
	Email    string      `json:"email" validate:"required,email" db:"email" example:"test@example.com"`
	Password string      `json:"password" validate:"required,min=8" u:"P1" db:"password" example:"12345678"`
	Check    int         `json:"check"  validate:"required,myvalidate"`
}
