package route

import "github.com/gin-gonic/gin"

func SetupRoutes(handlers *HandlersContainer) *gin.Engine {
	r := gin.Default()

	// Cart API
	r.POST("/cart", handlers.CartHandler.CreateCart)

	// Order API
	r.POST("/order", handlers.OrderHandler.CreateOrder)
	return r
}
