package handler

import (
	"fmt"
	"gotemplate/config"
	_ "gotemplate/docs"
	"sync/atomic"
	"github.com/gin-contrib/pprof"
	//"io"
	"net/http"
	//"os"
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

var isShuttingDown atomic.Value

func init() {
	isShuttingDown.Store(false)
}

// SetIsShuttingDown is an exported function that allows other packages to update the isShuttingDown value
func SetIsShuttingDown(shuttingDown bool) {
	isShuttingDown.Store(shuttingDown)
}

func HealthCheckHandler(c *gin.Context) {
	shuttingDown := isShuttingDown.Load().(bool)
	if shuttingDown {
		
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// NewRouter creates a new HTTP router
func NewRouter(
	cfg config.Econfig,
	userHandler UserHandler,
	paymentHandler PaymentHandler,
	categoryHandler CategoryHandler,
	productHandler ProductHandler,
	orderHandler OrderHandler,
	bagHander BagHandler,

) (*Router, error) {
	// Disable debug mode and write logs to file in production
	env := cfg.AppEnv()
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)

		//logFile, _ := os.Create("gin.log")
		//gin.DefaultWriter = io.Writer(logFile)
	}
	if env == "test" {
		gin.SetMode(gin.TestMode)
	}

	// CORS
	config := cors.DefaultConfig()

	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"*"}
	config.AllowBrowserExtensions = true
	config.AllowMethods = []string{"*"}

	router := gin.New()
	router.RedirectTrailingSlash = false
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": []string{"Invalid Path"},
			"errorno": []string{"INV1"},
		})

	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"success": false,
			"message": []string{"Method not Allowed"},
			"errorno": []string{"MD01"},
		})

	})

	//r.Use( ValidateContentType( []string{"application/json", "application/xml"}   )   )

	//router.Use(gin.LoggerWithFormatter(customLogger), gin.Recovery(), cors.New(config))




	// if env == "production" {
	// 	router.Use(gin.LoggerWithFormatter(customLogger), gin.Recovery(), cors.New(config), ValidateContentType([]string{"application/json"}))
	// }
	// if env == "development" {
	// 	router.Use(gin.LoggerWithFormatter(customLogger), gin.Recovery(), cors.New(config), ValidateContentType([]string{"application/json", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8", "text/css,*/*;q=0.1", "application/json,*/*", "*/*"}))

	// }





	
	//router.Use(gin.LoggerWithFormatter(customLogger), gin.Recovery(), cors.New(config))

	//router.Use(gin.LoggerWithFormatter(customLogger), gin.Recovery(), cors.New(config))

	/*	// Custom validators
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
				return nil, err
			}

			if err := v.RegisterValidation("payment_type", paymentTypeValidator); err != nil {
				return nil, err
			}

		}*/

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healthz", HealthCheckHandler)
	pprof.Register(router)
	v1 := router.Group("/v1")
	{ // @Router /users
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)

			authUser := user.Group("/")
			{
				authUser.GET("/", userHandler.ListUsers)
				authUser.GET("/:id", userHandler.GetUser)
				authUser.PUT("/:id", userHandler.UpdateUser)
				authUser.DELETE("/:id", userHandler.DeleteUser)

			}
		}

		bag := v1.Group("/bags")
		{
			bag.GET("/:id", bagHander.GetBag)
			bag.GET("/", bagHander.ListBags)
			bag.POST("/", bagHander.InsertBag)
			bag.POST("/sqrl", bagHander.Bagsquirrel)
			bag.POST("/pgx", bagHander.Bagspgx)
			bag.POST("/articles", bagHander.TxBagArticles)
			bag.POST("/all", bagHander.TxAllArrays)
			bag.POST("/insertpiece", bagHander.InsertPiece)
			bag.POST("/updatepiece", bagHander.updatepiece)
			//bag.POST("/updatepiecereturn", bagHander.updatepiecereturn)
			bag.POST("/updatepiecetx", bagHander.updatepiecetx)

		}

		payment := v1.Group("/payments")
		{
			payment.GET("/", paymentHandler.ListPayments)
			payment.GET("/:id", paymentHandler.GetPayment)
			payment.POST("/", paymentHandler.CreatePayment)
			payment.PUT("/:id", paymentHandler.UpdatePayment)
			payment.DELETE("/:id", paymentHandler.DeletePayment)

		}
		// @Router /categories
		category := v1.Group("/categories")
		{
			// @Router /v1/categories/ [get]
			category.GET("/", categoryHandler.ListCategories)
			category.GET("/:id", categoryHandler.GetCategory)
			category.POST("/", categoryHandler.CreateCategory)
			category.PUT("/:id", categoryHandler.UpdateCategory)
			category.DELETE("/:id", categoryHandler.DeleteCategory)

		}

		product := v1.Group("/products")
		{
			product.GET("/", productHandler.ListProducts)
			product.GET("/:id", productHandler.GetProduct)
			product.POST("/", productHandler.CreateProduct)
			product.PUT("/:id", productHandler.UpdateProduct)
			product.DELETE("/:id", productHandler.DeleteProduct)

		}
		order := v1.Group("/orders")
		{
			order.POST("/", orderHandler.CreateOrder)
			order.GET("/", orderHandler.ListOrders)
			order.GET("/:id", orderHandler.GetOrder)
		}
	}
	//}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}

func ValidateContentType(allowedTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		contentType := c.GetHeader("Content-Type")
		//var contentType string
		//contentType = c.GetHeader("Accept")
		//contentheader := c.ContentType()

		// if "Accept" == c.GetHeader("Accept") {
		// 	contentType = c.GetHeader("Accept")
		// }
		// if "Content-Type" == c.GetHeader("Content-Type") {
		// 	contentType = c.GetHeader("Content-Type")
		// }

		// Check if the Content-Type is in the allowedTypes
		validContentType := false
		for _, allowedType := range allowedTypes {
			if contentType == allowedType {
				validContentType = true
				break
			}
		}

		// if contentheader != "application/json" {
		// 	validContentType = false
		// }
		if !validContentType {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"success": false,
				"message": []string{"Invalid Content-Type. Supported types are: " + fmt.Sprintf("%v", allowedTypes)},
				"errorno": []string{"USP1"},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// customLogger is a custom Gin logger
func customLogger(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s [%s]\"\n",
		param.TimeStamp.Format(time.RFC1123),
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency.Round(time.Millisecond),
		param.Request.UserAgent(),
	)
}
