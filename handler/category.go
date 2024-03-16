package handler

import (
	"gotemplate/core/domain"
	"gotemplate/logger"

	repo "gotemplate/repo/postgres"

	"github.com/gin-gonic/gin"
)

// CategoryHandler represents the HTTP handler for category-related requests
type CategoryHandler struct {
	svc repo.CategoryRepository
	log *logger.Logger
	vs *ValidatorService
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(svc repo.CategoryRepository, log *logger.Logger,vs *ValidatorService) *CategoryHandler {
	return &CategoryHandler{
		svc,
		log,
		vs,
	}
}

// createCategoryRequest represents a request body for creating a new category
type createCategoryRequest struct {
	Name string `json:"name" validate:"required" example:"Foods"`
}

// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	create a new category with name
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			createCategoryRequest	body		createCategoryRequest	true	"Create category request"
//	@Success		200						{object}	categoryResponse		"Category created"
//	@Failure		400						{object}	errorValidResponse			"Validation error"
//	@Failure		401						{object}	errorValidResponse			"Unauthorized error"
//	@Failure		403						{object}	errorValidResponse			"Forbidden error"
//	@Failure		404						{object}	errorValidResponse			"Data not found error"
//	@Failure		409						{object}	errorValidResponse			"Data conflict error"
//	@Failure		500						{object}	errorValidResponse			"Internal server error"
//	@Router			/categories [post]

func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ch.vs.handleError(ctx, err)
		return
	}
	if !ch.vs.handleValidation(ctx, req) {
		return
	}
	category := domain.Category{
		Name: req.Name,
	}

	_, err := ch.svc.CreateCategory(ctx, &category)
	if err != nil {
		ch.log.Error(err.Error())
		ch.vs.handledbError(ctx, err)
		return
	}

	rsp := newCategoryResponse(&category)

	handleSuccess(ctx, rsp)
}

// getCategoryRequest represents a request body for retrieving a category
type getCategoryRequest struct {
	ID uint64 `uri:"id" validate:"required,min=1" example:"1"`
}

// GetCategory godoc
//
//	@Summary		Get a category
//	@Description	get a category by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64				true	"Category ID"
//	@Success		200	{object}	categoryResponse	"Category retrieved"
//	@Failure		400	{object}	errorValidResponse		"Validation error"
//	@Failure		404	{object}	errorValidResponse		"Data not found error"
//	@Failure		500	{object}	errorValidResponse		"Internal server error"
//	@Router			/categories/{id} [get]
func (ch *CategoryHandler) GetCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ch.vs.handleError(ctx, err)
		return
	}
	if !ch.vs.handleValidation(ctx, req) {
		return
	}
	category, err := ch.svc.GetCategoryByID(ctx, req.ID)
	if err != nil {
		ch.log.Error(err.Error())
		ch.vs.handledbError(ctx, err)
		return
	}

	rsp := newCategoryResponse(category)

	handleSuccess(ctx, rsp)
}

// listCategoriesRequest represents a request body for listing categories
type listCategoriesRequest struct {
	Skip  uint64 `form:"skip" validate:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" validate:"required,min=5" example:"5"`
}

// ListCategories godoc
//
//	@Summary		List categories
//	@Description	List categories with pagination
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	categoryResponse			"Categories displayed"
//	@Failure		400		{object}	errorValidResponse	"Validation error"
//	@Failure		500		{object}	errorValidResponse	"Internal server error"
//	@Router			/categories/ [get]
func (ch *CategoryHandler) ListCategories(ctx *gin.Context) {
	var req listCategoriesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ch.vs.handleError(ctx, err)
		return
	}
	if !ch.vs.handleValidation(ctx, req) {
		return
	}
	categories, err := ch.svc.ListCategories(ctx, req.Skip, req.Limit)
	if err != nil {
		ch.log.Error(err.Error())
		ch.vs.handledbError(ctx, err)
		return
	}

	total := uint64(len(categories))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, categories, "categories")

	handleSuccess(ctx, rsp)
}

// updateCategoryRequest represents a request body for updating a category
type updateCategoryRequest struct {
	Name string `json:"name" validate:"omitempty,required" example:"Beverages"`
}

// UpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	update a category's name by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id						path		uint64					true	"Category ID"
//	@Param			updateCategoryRequest	body		updateCategoryRequest	true	"Update category request"
//	@Success		200						{object}	categoryResponse		"Category updated"
//	@Failure		400						{object}	errorValidResponse			"Validation error"
//	@Failure		401						{object}	errorValidResponse			"Unauthorized error"
//	@Failure		403						{object}	errorValidResponse			"Forbidden error"
//	@Failure		404						{object}	errorValidResponse			"Data not found error"
//	@Failure		409						{object}	errorValidResponse			"Data conflict error"
//	@Failure		500						{object}	errorValidResponse			"Internal server error"
//	@Router			/categories/{id} [put]
func (ch *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var req updateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ch.vs.handleError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		ch.vs.handleError(ctx, err)
		return
	}
	if !ch.vs.handleValidation(ctx, req) {
		return
	}
	category := domain.Category{
		ID:   id,
		Name: req.Name,
	}

	_, err = ch.svc.UpdateCategory(ctx, &category)
	if err != nil {
		ch.log.Error(err.Error())
		ch.vs.handledbError(ctx, err)
		return
	}

	rsp := newCategoryResponse(&category)

	handleSuccess(ctx, rsp)
}

// deleteCategoryRequest represents a request body for deleting a category
type deleteCategoryRequest struct {
	ID uint64 `uri:"id" validate:"required,min=1" example:"1"`
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Delete a category by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Category ID"
//	@Success		200	{object}	categoryResponse		"Category deleted"
//	@Failure		400	{object}	errorValidResponse	"Validation error"
//	@Failure		401	{object}	errorValidResponse	"Unauthorized error"
//	@Failure		403	{object}	errorValidResponse	"Forbidden error"
//	@Failure		404	{object}	errorValidResponse	"Data not found error"
//	@Failure		500	{object}	errorValidResponse	"Internal server error"
//	@Router			/categories/{id} [delete]
func (ch *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ch.vs.handleError(ctx, err)
		return
	}
	if !ch.vs.handleValidation(ctx, req) {
		return
	}
	err := ch.svc.DeleteCategory(ctx, req.ID)
	if err != nil {
		ch.log.Error(err.Error())
		ch.vs.handledbError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
