package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter(

	userHandler UserHandler,
	paymentHandler PaymentHandler,
	categoryHandler CategoryHandler,
	productHandler ProductHandler,
	orderHandler OrderHandler,

) (*Router, error) {
	// Disable debug mode and write logs to file in production
	env := os.Getenv("APP_ENV")
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
	router.Use(gin.LoggerWithFormatter(customLogger), gin.Recovery(), cors.New(config), ValidateContentType([]string{"application/json"}))

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
	//router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
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
		payment := v1.Group("/payments")
		{
			payment.GET("/", paymentHandler.ListPayments)
			payment.GET("/:id", paymentHandler.GetPayment)
			payment.POST("/", paymentHandler.CreatePayment)
			payment.PUT("/:id", paymentHandler.UpdatePayment)
			payment.DELETE("/:id", paymentHandler.DeletePayment)

		}

		category := v1.Group("/categories")
		{
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

		// Check if the Content-Type is in the allowedTypes
		validContentType := false
		for _, allowedType := range allowedTypes {
			if contentType == allowedType {
				validContentType = true
				break
			}
		}

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
