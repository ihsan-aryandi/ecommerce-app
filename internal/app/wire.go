//go:build wireinject
// +build wireinject

package app

import (
	"ecommerce-app/internal/config"
	"ecommerce-app/internal/database"
	"ecommerce-app/internal/provider"
	"ecommerce-app/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeApp() *gin.Engine {
	wire.Build(
		// Config
		config.Load,

		// Database
		database.CreateConnection,

		// Handlers
		provider.HandlersSet,

		// Services
		provider.ServicesSet,

		// Repositories
		provider.RepositoriesSet,

		// Handlers Container
		route.NewHandlersContainer,

		// Routes Setup
		route.SetupRoutes,
	)

	return nil
}
