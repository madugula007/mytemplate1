package handler

import (
	//"database/sql"

	"gotemplate/core/domain"

	//"github.com/templatedop/githubrepo/dtime"
	"gotemplate/logger"

	//"time"

	//"github.com/guregu/null"
	"github.com/guregu/null/zero"

	//"github.com/jackc/pgx/v5/pgtype"
	//"github.com/aarondl/opt/null"
	"github.com/jinzhu/copier"

	//"github.com/volatiletech/null"

	//"gotemplate/core/port"
	repo "gotemplate/repo/postgres"

	"github.com/gin-gonic/gin"
	//"gotemplate/dtime"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc repo.UserRepository
	log *logger.Logger
	vs  *ValidatorService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc repo.UserRepository, log *logger.Logger, vs *ValidatorService) *UserHandler {
	return &UserHandler{
		svc,
		log,
		vs,
	}
}

// registerRequest represents the request body for creating a user
type RegisterRequest struct {
	//Name     string `json:"name" validate:"omitempty,required,min=5" u:"N1" db:"name" example:"John Doe"`
	Name        string `json:"name" validate:"omitempty,min=5" u:"N1" db:"name" example:"John Doe"`
	Email       string `json:"email" validate:"required,email" db:"email" example:"test@example.com"`
	Password    string `json:"password" validate:"required,min=8" u:"P1" db:"password" example:"12345678"`
	Check       int    `json:"check" validate:"required,myvalidate"`
	CreatedAt   string `json:"created_at" `
	CreatedTime string ` json:"created_time"  db:"created_time"  validate:"required,hourvalidate"`
}

type User struct{}

func (uh *UserHandler) Register(ctx *gin.Context) {

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		uh.log.Debug("Came inside")
		uh.vs.handleError(ctx, err)
		return
	}

	if !uh.vs.handleValidation(ctx, req) {
		return
	}

	//uh.log.Debug("req:", req)

	user := domain.User{

		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		//CreatedAt: null.StringFrom(req.CreatedAt),
		//CreatedAt: null.From(req.CreatedAt),

		//null.StringFrom()

		CreatedAt: zero.StringFrom(req.CreatedAt),
		//CreatedAt: req.CreatedAt,
		//CreatedTime: dtime.NewFromStrFormat(req.CreatedTime, "h:i"),
		//Role:     req.Role,
	}

	//var user domain.User

	// err := copier.Copy(&user, req)
	// if err != nil {
	// 	uh.log.Debug("error: ", err)
	// }

	//uh.log.Debug("User:", user)
	u, err := uh.svc.CreateUser(ctx, &user)
	if err != nil {
		uh.log.Error(err.Error())
		uh.vs.handledbError(ctx, err)
		return
	}
	// u1 := UserResponse{}
	// copier.Copy(&u1, &u)
	// rsp := newUserResponse(&u1)

	// u, err := uh.svc.CreateUser(ctx, &user)
	// if err != nil {
	// 	uh.log.Error(err.Error())
	// 	handledbError(ctx, err)
	// 	return
	// }
	// u1 := UserResponse{}
	// copier.Copy(&u1, &u)
	// rsp := newUserResponse(&u1)

	handleSuccess(ctx, u)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Skip  uint64 `form:"skip" validate:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" validate:"required,min=5" example:"5"`
}

func (uh *UserHandler) ListUsers(ctx *gin.Context) {

	var req listUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		uh.vs.handleError(ctx, err)
		return
	}
	if !uh.vs.handleValidation(ctx, req) {
		return
	}

	users, err := uh.svc.ListUsers(ctx, req.Skip, req.Limit)
	if err != nil {
		uh.log.Error(err.Error())
		uh.vs.handledbError(ctx, err)
		return
	}

	total := uint64(len(users))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, users, "users")

	handleSuccess(ctx, rsp)
}

// getUserRequest represents the request body for getting a user
type getUserRequest struct {
	ID uint64 `uri:"id" validate:"required,min=1" example:"1"`
}

func (uh *UserHandler) GetUser(ctx *gin.Context) {

	lengthStr := ctx.Param("Length")
	uh.log.Debug("lenght string", lengthStr)
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		uh.vs.handleError(ctx, err)
		return
	}
	if !uh.vs.handleValidation(ctx, req) {
		return
	}
	user, b, err := uh.svc.GetUserByID(ctx, req.ID)
	if err != nil {

		uh.log.Error(err.Error())
		uh.vs.handledbError(ctx, err)
		return
	}

	if b {

		uh.log.Debug("User before copying:", user)
		// userr := UserResponse{
		// 	ID:    user.ID,
		// 	Name:  user.Name,
		// 	Email: user.Email,
		// 	//CreatedTime: user.CreatedTime.String(),
		// 	//CreatedTime: user.CreatedTime,
		// 	//Role:     req.Role,
		// }

		u := UserResponse{}
		copier.Copy(&u, user)
		// uh.log.Debug("user:", u)

		rsp := newUserResponse1(u)

		handleSuccess(ctx, rsp)
	} else {

		handleSuccess(ctx, "No Rows")
	}

}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Name     string          `json:"name" validate:"omitempty,required" example:"John Doe"`
	Email    string          `json:"email" validate:"omitempty,required,email" example:"test@example.com"`
	Password string          `json:"password" validate:"omitempty,required,min=8" example:"12345678"`
	Role     domain.UserRole `json:"role" validate:"omitempty,required" example:"admin"`
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		uh.vs.handleError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		uh.vs.handleError(ctx, err)
		return
	}

	if !uh.vs.handleValidation(ctx, req) {
		return
	}

	user := domain.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		//Role:     req.Role,
	}

	_, err = uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		uh.log.Error(err.Error())
		uh.vs.handledbError(ctx, err)
		return
	}
	u := UserResponse{}
	copier.Copy(&u, &user)
	rsp := newUserResponse(&u)

	handleSuccess(ctx, rsp)
}

// deleteUserRequest represents the request body for deleting a user
type deleteUserRequest struct {
	ID uint64 `uri:"id" validate:"required,min=1" example:"1"`
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		uh.vs.handleError(ctx, err)
		return
	}
	if !uh.vs.handleValidation(ctx, req) {
		return
	}
	err := uh.svc.DeleteUser(ctx, req.ID)
	if err != nil {
		uh.log.Error(err.Error())
		uh.vs.handledbError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
