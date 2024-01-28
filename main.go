package main

import (
	"context"
	config "gotemplate/config"

	//handler "gotemplate/handler"
	"gotemplate/logger"
	repo "gotemplate/repo/postgres"
	r "gotemplate/route"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	config.Load()
}

func LoggerInit(loglevel string) *logger.Logger {
	log := logger.New(loglevel)
	return log
}

func main() {

	log := LoggerInit(os.Getenv("LOG_LEVEL"))

	log.Debug("Hi... Am I visible")
	log.Info("This is Info")
	log.Error("This is error")
	log.Warn("This is warning")
	listenAddr := ":" + os.Getenv("HTTP_PORT")
	ctx := context.Background()
	db, err := repo.NewDB(ctx)
	if err != nil {
		log.Warn("error in db connection %s", err)
	}
	defer db.Close()
	log.Info("Successfully connected to the database %s", os.Getenv("DB_CONNECTION"))

	router, err := r.Routes(db, log)

	/*userRepo := repo.NewUserRepository(db, log)

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



	router, err := handler.NewRouter(
		*userHandler,
		*paymentHandler,
		*categoryHandler,
		*productHandler,
		*orderHandler,
	)*/
	if err != nil {
		log.Warn("Error initializing router %s", err)
		os.Exit(1)

	}

	log.Info("Starting the HTTP server: %s", listenAddr)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("HTTP_PORT"),
		Handler: router,
	}

	go func() {

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("Received signal...%s", sig)

	duration, err := time.ParseDuration(os.Getenv("SHUTDOWN_TIME"))
	if err != nil {
		log.Fatal("Error in parsing duration", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown error:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Info("Server exiting")

}
