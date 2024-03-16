package route

import (
	"gotemplate/config"
	handler "gotemplate/handler"
	"gotemplate/logger"
	repo "gotemplate/repo/postgres"
)

//	func LoggerInit(loglevel string) *logger.Logger {
//		log := logger.New(loglevel)
//		return log
//	}
func Routes(db *repo.DB, log *logger.Logger, cfg config.Econfig, validatorService *handler.ValidatorService) (router *handler.Router, err1 error) {

	userRepo := repo.NewUserRepository(db, log)

	userHandler := handler.NewUserHandler(*userRepo, log, validatorService)

	bagRepo := repo.NewBagRepository(db, log)

	bagHandler := handler.NewBagHandler(*bagRepo, log, validatorService)

	paymentRepo := repo.NewPaymentRepository(db, log)
	paymentHandler := handler.NewPaymentHandler(*paymentRepo, log, validatorService)

	// Category
	categoryRepo := repo.NewCategoryRepository(db, log)
	categoryHandler := handler.NewCategoryHandler(*categoryRepo, log, validatorService)

	// Product
	productRepo := repo.NewProductRepository(db, log)
	productHandler := handler.NewProductHandler(*productRepo, log, validatorService)

	// Order
	orderRepo := repo.NewOrderRepository(db, log)
	orderHandler := handler.NewOrderHandler(*orderRepo, log, validatorService)

	router, err1 = handler.NewRouter(
		cfg,
		*userHandler,
		*paymentHandler,
		*categoryHandler,
		*productHandler,
		*orderHandler,
		*bagHandler,
	)
	return router, err1

}
