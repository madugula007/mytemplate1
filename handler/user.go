package handler

import (
	"gotemplate/core/domain"
	"gotemplate/logger"

	"github.com/jinzhu/copier"

	//"gotemplate/core/port"
	repo "gotemplate/repo/postgres"

	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc repo.UserRepository
	log *logger.Logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc repo.UserRepository, log *logger.Logger) *UserHandler {
	return &UserHandler{
		svc,
		log,
	}
}

// registerRequest represents the request body for creating a user
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=5" u:"N1" db:"name" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" db:"email" example:"test@example.com"`
	Password string `json:"password" validate:"required,min=8" u:"P1" db:"password" example:"12345678"`
	Check    int    `json:"check" validate:"required,myvalidate"`
}

type User struct{}

func (uh *UserHandler) Register(ctx *gin.Context) {

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		uh.log.Debug("Came inside")
		handleError(ctx, err)
		return
	}

	if !handleValidation(ctx, req) {
		return
	}

	var user domain.User
	// user := domain.User{
	// 	Name:     req.Name,
	// 	Email:    req.Email,
	// 	Password: req.Password,
	// }

	err := copier.Copy(&user, &req)
	if err != nil {
		uh.log.Debug("error: ", err)
	}

	// uh.log.Debug("name:", user.Name)
	// uh.log.Debug("email:", user.Email)
	// uh.log.Debug("password:", user.Password)

	u, err := uh.svc.CreateUser(ctx, &user)
	if err != nil {
		uh.log.Error(err.Error())
		handledbError(ctx, err)
		return
	}
	u1 := UserResponse{}
	copier.Copy(&u1, &u)
	//rsp := newUserResponse(u)

	handleSuccess(ctx, u1)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

func (uh *UserHandler) ListUsers(ctx *gin.Context) {

	var req listUsersRequest
	//var usersList []userResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	users, err := uh.svc.ListUsers(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// for _, user := range users {
	// 	usersList = append(usersList, newUserResponse(&user))
	// }

	total := uint64(len(users))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, users, "users")

	handleSuccess(ctx, rsp)
}

// getUserRequest represents the request body for getting a user
type getUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

func (uh *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, err := uh.svc.GetUserByID(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(user)

	handleSuccess(ctx, rsp)
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Name     string          `json:"name" binding:"omitempty,required" example:"John Doe"`
	Email    string          `json:"email" binding:"omitempty,required,email" example:"test@example.com"`
	Password string          `json:"password" binding:"omitempty,required,min=8" example:"12345678"`
	Role     domain.UserRole `json:"role" binding:"omitempty,required" example:"admin"`
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	_, err = uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(&user)

	handleSuccess(ctx, rsp)
}

// deleteUserRequest represents the request body for deleting a user
type deleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := uh.svc.DeleteUser(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
