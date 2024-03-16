package handler

import (
	"gotemplate/core/domain"
	"gotemplate/logger"

	"github.com/gin-gonic/gin"

	repo "gotemplate/repo/postgres"
)

// ProductHandler represents the HTTP handler for product-related requests
type ProductHandler struct {
	svc repo.ProductRepository
	log *logger.Logger
	vs *ValidatorService
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(svc repo.ProductRepository, log *logger.Logger,vs *ValidatorService) *ProductHandler {
	return &ProductHandler{
		svc,
		log,
		vs,
	}
}

// createProductRequest represents a request body for creating a new product
type createProductRequest struct {
	CategoryID uint64  `json:"category_id" validate:"required,min=1" example:"1"`
	Name       string  `json:"name" validate:"required" example:"Chiki Ball"`
	Image      string  `json:"image" validate:"required" example:"https://example.com/chiki-ball.png"`
	Price      float64 `json:"price" validate:"required,min=0" example:"5000"`
	Stock      int64   `json:"stock" validate:"required,min=0" example:"100"`
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	create a new product with name, image, price, and stock
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			createProductRequest	body		createProductRequest	true	"Create product request"
//	@Success		200						{object}	productResponse			"Product created"
//	@Failure		400						{object}	errorValidResponse			"Validation error"
//	@Failure		401						{object}	errorValidResponse			"Unauthorized error"
//	@Failure		403						{object}	errorValidResponse			"Forbidden error"
//	@Failure		404						{object}	errorValidResponse			"Data not found error"
//	@Failure		409						{object}	errorValidResponse			"Data conflict error"
//	@Failure		500						{object}	errorValidResponse			"Internal server error"
//	@Router			/products [post]
func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ph.vs.handleError(ctx, err)
		return
	}
	if !ph.vs.handleValidation(ctx, req) {
		return
	}

	product := domain.Product{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Image:      req.Image,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	_, err := ph.svc.CreateProduct(ctx, &product)
	if err != nil {
		ph.log.Error(err.Error())
		ph.vs.handledbError(ctx, err)
		return
	}

	rsp := newProductResponse(&product)

	handleSuccess(ctx, rsp)
}

// getProductRequest represents a request body for retrieving a product
type getProductRequest struct {
	ID uint64 `uri:"id" validate:"required,min=1" example:"1"`
}

// GetProduct godoc
//
//	@Summary		Get a product
//	@Description	get a product by id with its category
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Product ID"
//	@Success		200	{object}	productResponse	"Product retrieved"
//	@Failure		400	{object}	errorValidResponse	"Validation error"
//	@Failure		404	{object}	errorValidResponse	"Data not found error"
//	@Failure		500	{object}	errorValidResponse	"Internal server error"
//	@Router			/products/{id} [get]
func (ph *ProductHandler) GetProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ph.vs.handleError(ctx, err)
		return
	}
	if !ph.vs.handleValidation(ctx, req) {
		return
	}
	product, err := ph.svc.GetProductByID(ctx, req.ID)
	if err != nil {
		ph.log.Error(err.Error())
		ph.vs.handledbError(ctx, err)
		return
	}

	rsp := newProductResponse(product)

	handleSuccess(ctx, rsp)
}

// listProductsRequest represents a request body for listing products
type listProductsRequest struct {
	CategoryID uint64 `form:"category_id" validate:"omitempty,min=1" example:"1"`
	Query      string `form:"q" validate:"omitempty" example:"Chiki"`
	Skip       uint64 `form:"skip" validate:"required,min=0" example:"0"`
	Limit      uint64 `form:"limit" validate:"required,min=5" example:"5"`
}

// ListProducts godoc
//
//	@Summary		List products
//	@Description	List products with pagination
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			category_id	query		uint64			false	"Category ID"
//	@Param			q			query		string			false	"Query"
//	@Param			skip		query		uint64			true	"Skip"
//	@Param			limit		query		uint64			true	"Limit"
//	@Success		200			{object}	meta			"Products retrieved"
//	@Failure		400			{object}	errorValidResponse	"Validation error"
//	@Failure		500			{object}	errorValidResponse	"Internal server error"
//	@Router			/products [get]
func (ph *ProductHandler) ListProducts(ctx *gin.Context) {
	var req listProductsRequest
	//var productsList []productResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ph.vs.handleError(ctx, err)
		return
	}
	if !ph.vs.handleValidation(ctx, req) {
		return
	}
	products, err := ph.svc.ListProducts(ctx, req.Query, req.CategoryID, req.Skip, req.Limit)
	if err != nil {
		ph.log.Error(err.Error())
		ph.vs.handledbError(ctx, err)
		return
	}

	total := uint64(len(products))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, products, "products")

	handleSuccess(ctx, rsp)
}

// updateProductRequest represents a request body for updating a product
type updateProductRequest struct {
	CategoryID uint64  `json:"category_id" validate:"omitempty,required,min=1" example:"1"`
	Name       string  `json:"name" validate:"omitempty,required" example:"Nutrisari Jeruk"`
	Image      string  `json:"image" validate:"omitempty,required" example:"https://example.com/nutrisari-jeruk.png"`
	Price      float64 `json:"price" validate:"omitempty,required,min=0" example:"2000"`
	Stock      int64   `json:"stock" validate:"omitempty,required,min=0" example:"200"`
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	update a product's name, image, price, or stock by id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id						path		uint64					true	"Product ID"
//	@Param			updateProductRequest	body		updateProductRequest	true	"Update product request"
//	@Success		200						{object}	productResponse			"Product updated"
//	@Failure		400						{object}	errorValidResponse			"Validation error"
//	@Failure		401						{object}	errorValidResponse			"Unauthorized error"
//	@Failure		403						{object}	errorValidResponse			"Forbidden error"
//	@Failure		404						{object}	errorValidResponse			"Data not found error"
//	@Failure		409						{object}	errorValidResponse			"Data conflict error"
//	@Failure		500						{object}	errorValidResponse			"Internal server error"
//	@Router			/products/{id} [put]
func (ph *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ph.vs.handleError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		ph.vs.handleError(ctx, err)
		return
	}
	if !ph.vs.handleValidation(ctx, req) {
		return
	}
	product := domain.Product{
		ID:         id,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Image:      req.Image,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	_, err = ph.svc.UpdateProduct(ctx, &product)
	if err != nil {
		ph.log.Error(err.Error())
		ph.vs.handledbError(ctx, err)
		return
	}

	rsp := newProductResponse(&product)

	handleSuccess(ctx, rsp)
}

// deleteProductRequest represents a request body for deleting a product
type deleteProductRequest struct {
	ID uint64 `uri:"id" validate:"required,min=1" example:"1"`
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product by id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Product ID"
//	@Success		200	{object}	categoryResponse		"Product deleted"
//	@Failure		400	{object}	errorValidResponse	"Validation error"
//	@Failure		401	{object}	errorValidResponse	"Unauthorized error"
//	@Failure		403	{object}	errorValidResponse	"Forbidden error"
//	@Failure		404	{object}	errorValidResponse	"Data not found error"
//	@Failure		500	{object}	errorValidResponse	"Internal server error"
//	@Router			/products/{id} [delete]

func (ph *ProductHandler) DeleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ph.vs.handleError(ctx, err)
		return
	}
	if !ph.vs.handleValidation(ctx, req) {
		return
	}
	err := ph.svc.DeleteProduct(ctx, req.ID)
	if err != nil {
		ph.log.Error(err.Error())
		ph.vs.handledbError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
