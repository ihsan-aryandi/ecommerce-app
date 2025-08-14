package route

import "github.com/gin-gonic/gin"

func SetupRoutes(handlers *HandlersContainer) *gin.Engine {
	r := gin.Default()

	// Checkout API
	r.GET("/checkout", handlers.CheckoutHandler.GetTotal)
	return r
}
