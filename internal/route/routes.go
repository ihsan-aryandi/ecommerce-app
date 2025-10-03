package route

import "github.com/gin-gonic/gin"

func SetupRoutes(handlers *HandlersContainer) *gin.Engine {
	r := gin.Default()

	// Cart API
	r.POST("/cart", handlers.CartHandler.CreateCart)

	// Checkout API
	r.POST("/checkout", handlers.CheckoutHandler.Checkout)
	r.POST("/checkout/summary", handlers.CheckoutHandler.CalculateSummaries)

	// Order API
	r.POST("/order", handlers.OrderHandler.CreateOrder)
	return r
}
