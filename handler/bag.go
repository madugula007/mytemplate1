package handler

import (
	"gotemplate/core/domain"
	"gotemplate/logger"

	//"gotemplate/core/port"
	repo "gotemplate/repo/postgres"

	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for user-related requests
type BagHandler struct {
	svc repo.BagRepository
	log *logger.Logger
	vs  *ValidatorService
}

// NewUserHandler creates a new UserHandler instance
func NewBagHandler(svc repo.BagRepository, log *logger.Logger, vs *ValidatorService) *BagHandler {
	return &BagHandler{
		svc,
		log,
		vs,
	}
}

func (bh *BagHandler) InsertPiece(ctx *gin.Context) {

	var req domain.InternationalArticleSubpiece
	//bh.log.Debug("req:",req)
	if err := ctx.BindJSON(&req); err != nil {
		bh.vs.handleError(ctx, err)
		return
	}
	if !bh.vs.handleValidation(ctx, req) {
		return
	}
	err := bh.svc.InsertPiece(ctx, &req)
	if err != nil {
		bh.log.Error(err.Error())
		bh.vs.handledbError(ctx, err)
		return
	}
	// u := UserResponse{}
	// copier.Copy(&u, user)

	//rsp := newUserResponse1(bag)

	handleSuccess(ctx, req)
}

func (bh *BagHandler) InsertBag(ctx *gin.Context) {
	var req domain.Bag1
	//bh.log.Debug("req:",req)
	if err := ctx.BindJSON(&req); err != nil {
		bh.vs.handleError(ctx, err)
		return
	}
	if !bh.vs.handleValidation(ctx, req) {
		return
	}
	bag, err := bh.svc.Insertbag(ctx, req)
	if err != nil {
		bh.log.Error(err.Error())
		bh.vs.handledbError(ctx, err)
		return
	}
	// u := UserResponse{}
	// copier.Copy(&u, user)

	//rsp := newUserResponse1(bag)

	handleSuccess(ctx, bag)
}

type combinedstruct struct {
	subp  domain.InternationalArticleSubpiece
	asubp []domain.InternationalArticleSubpiece
}

// func (bh *BagHandler) updatepiecereturn(ctx *gin.Context) {
// 	var req domain.ISubpieces
// 	//bh.log.Debug("req:",req)
// 	if err := ctx.BindJSON(&req); err != nil {
// 		bh.vs.handleError(ctx, err)
// 		return
// 	}
// 	if !bh.vs.handleValidation(ctx, req) {
// 		return
// 	}

// 	//bh.log.Debug("req:",req.IntlSubpieces.)

// 	b, err := bh.svc.Updatepieceswithreturn(ctx, req)
// 	if err != nil {
// 		bh.log.Error(err.Error())
// 		bh.vs.handledbError(ctx, err)
// 		return
// 	}
// 	// u := UserResponse{}
// 	// copier.Copy(&u, user)

// 	//rsp := newUserResponse1(bag)

// 	handleSuccess(ctx, b)
// }

func (bh *BagHandler) updatepiecetx(ctx *gin.Context) {
	//bh.log.Debug("Came here:")
	var req domain.ISubpieces
	//bh.log.Debug("req:",req)
	if err := ctx.BindJSON(&req); err != nil {
		bh.vs.handleError(ctx, err)
		return
	}
	if !bh.vs.handleValidation(ctx, req) {
		return
	}

	//bh.log.Debug("req:",req.IntlSubpieces.)

	//bh.log.Debug("len of sub pieces:", len(req.IntlSubpieces))

	a, err := bh.svc.UpdatepieceswithTransaction(ctx, req)
	if err != nil {
		bh.log.Error(err.Error())
		bh.vs.handledbError(ctx, err)
		return
	}

	handleSuccess(ctx, a)
}

func (bh *BagHandler) updatepiece(ctx *gin.Context) {
	//bh.log.Debug("Came here:")
	var req domain.ISubpieces
	//bh.log.Debug("req:",req)
	if err := ctx.BindJSON(&req); err != nil {
		bh.vs.handleError(ctx, err)
		return
	}
	if !bh.vs.handleValidation(ctx, req) {
		return
	}
	a, err := bh.svc.Updatepieceswithbatch(ctx, req)
	if err != nil {
		bh.log.Error(err.Error())
		bh.vs.handledbError(ctx, err)
		return
	}

	handleSuccess(ctx, a)
}

func (bh *BagHandler) GetBag(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		bh.vs.handleError(ctx, err)
		return
	}
	if !bh.vs.handleValidation(ctx, req) {
		return
	}
	bag, err := bh.svc.GetBagByID(ctx, req.ID)
	if err != nil {
		bh.log.Error(err.Error())
		bh.vs.handledbError(ctx, err)
		return
	}
	// u := UserResponse{}
	// copier.Copy(&u, user)

	//rsp := newUserResponse1(bag)

	handleSuccess(ctx, bag)
}

type listBagsRequest struct {
	Skip  uint64 `form:"skip" validate:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" validate:"required,min=5" example:"5"`
}

func (ub *BagHandler) ListBags(ctx *gin.Context) {

	var req listBagsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ub.vs.handleError(ctx, err)
		return
	}
	if !ub.vs.handleValidation(ctx, req) {
		return
	}

	bags, err := ub.svc.GetBags(ctx, req.Skip, req.Limit)
	if err != nil {
		ub.log.Error(err.Error())
		ub.vs.handledbError(ctx, err)
		return
	}

	total := uint64(len(bags))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, bags, "bags")
	handleSuccess(ctx, rsp)
}

func (ub *BagHandler) Bagsquirrel(ctx *gin.Context) {
	var bags domain.Bags
	if err := ctx.ShouldBindJSON(&bags); err != nil {
		ub.log.Error("error at should bind json", err.Error())
		ub.vs.handleError(ctx, err)
		return
	}
	if !ub.vs.handleValidation(ctx, bags) {
		return
	}
	err := ub.svc.Insertbagswithsquirrel(ctx, bags)
	if err != nil {
		ub.log.Error(err.Error())
		ub.vs.handledbError(ctx, err)
		return
	}
	handleSuccess(ctx, "")
}

func (ub *BagHandler) Bagspgx(ctx *gin.Context) {
	var bags domain.Bags
	if err := ctx.ShouldBindJSON(&bags); err != nil {
		ub.vs.handleError(ctx, err)
		return
	}
	if !ub.vs.handleValidation(ctx, bags) {
		return
	}
	err := ub.svc.Insertbagswithpgx(ctx, bags)
	if err != nil {
		ub.log.Error(err.Error())
		ub.vs.handledbError(ctx, err)
		return
	}
	handleSuccess(ctx, "")
}

type bagarticles struct {
	domain.Article `json:"articles"`
	domain.Bag     `json:"bags"`
}

func (ub *BagHandler) TxBagArticles(ctx *gin.Context) {
	var bagarts bagarticles

	//how to separate the articles think
	if err := ctx.ShouldBindJSON(&bagarts); err != nil {
		ub.vs.handleError(ctx, err)
		return
	}
	if !ub.vs.handleValidation(ctx, bagarts) {
		return
	}

	err := repo.Tx(ctx, ub.svc.Db, ub.svc.InsertBagArticle, bagarts.Bag, bagarts.Article)
	//ub.svc.InsertBagArticle(ctx, bag, article)
	if err != nil {
		ub.log.Error(err.Error())
		ub.vs.handledbError(ctx, err)
		return
	}
	handleSuccess(ctx, "")
}

//func (br *BagRepository) InsertDataBulk(ctx context.Context, bags []domain.Bag, articles []domain.Article, phones []domain.Phone) error {

type alltxns struct {
	Bags     []domain.Bag1    `json:"bags"`
	Articles []domain.Article `json:"articles"`
	Phones   []domain.Phone   `json:"phones"`
}

func (ub *BagHandler) TxAllArrays(ctx *gin.Context) {

	var bag alltxns
	//how to separate the articles think
	if err := ctx.ShouldBindJSON(&bag); err != nil {
		ub.vs.handleError(ctx, err)
		return
	}
	if !ub.vs.handleValidation(ctx, bag) {
		return
	}

	err := repo.Tx(ctx, ub.svc.Db, ub.svc.InsertDataBulk, bag.Articles, bag.Bags, bag.Phones)
	//err := ub.svc.InsertDataBulk(ctx, bag, article, phone)
	if err != nil {
		ub.log.Error(err.Error())
		ub.vs.handledbError(ctx, err)
		return
	}
	handleSuccess(ctx, "")
}
