package route

import "github.com/gin-gonic/gin"

func SetupRoutes(handlers *HandlersContainer) *gin.Engine {
	r := gin.Default()

	// Cart API
	r.POST("/cart", handlers.CartHandler.CreateCart)

	// Calculate API
	r.POST("/calculate-summary", handlers.CalculateHandler.CalculateSummaries)

	// Order API
	r.POST("/order", handlers.OrderHandler.CreateOrder)
	return r
}
