package route

import (
	handler "gotemplate/handler"
	"gotemplate/logger"
	repo "gotemplate/repo/postgres"
)

func LoggerInit(loglevel string) *logger.Logger {
	log := logger.New(loglevel)
	return log
}
func Routes(db *repo.DB, log *logger.Logger) (router *handler.Router, err1 error) {

	// User

	userRepo := repo.NewUserRepository(db, log)

	userHandler := handler.NewUserHandler(*userRepo, log)

	paymentRepo := repo.NewPaymentRepository(db)
	paymentHandler := handler.NewPaymentHandler(*paymentRepo)

	// Category
	categoryRepo := repo.NewCategoryRepository(db)
	categoryHandler := handler.NewCategoryHandler(*categoryRepo)

	// Product
	productRepo := repo.NewProductRepository(db)
	productHandler := handler.NewProductHandler(*productRepo)

	// Order
	orderRepo := repo.NewOrderRepository(db)
	orderHandler := handler.NewOrderHandler(*orderRepo)

	router, err1 = handler.NewRouter(
		*userHandler,
		*paymentHandler,
		*categoryHandler,
		*productHandler,
		*orderHandler,
	)
	return router, err1

}
