package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malikkhoiri/csms/internal/domain"
)

type APIHandler struct {
	port string
}

func NewAPIHandler(port string) *APIHandler {
	return &APIHandler{
		port: port,
	}
}

func (h *APIHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
		"port":   h.port,
	})
}

func (h *APIHandler) GetConnections(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connections": "WebSocket connections info will be implemented",
	})
}

func SetupRoutes(
	router *gin.Engine,
	chargePointService domain.ChargePointService,
	transactionService domain.TransactionService,
	userService domain.UserService,
	idTagService domain.IDTagService,
	authService domain.AuthService,
) {
	authHandler := NewAuthHandler(authService)
	dashboardHandler := NewDashboardHandler(chargePointService, transactionService, userService)
	chargePointHandler := NewChargePointHandler(chargePointService)
	transactionHandler := NewTransactionHandler(transactionService)
	userHandler := NewUserHandler(userService)
	idTagHandler := NewIDTagHandler(idTagService)

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.Me)
	}

	api := router.Group("/api/v1")
	api.Use(AuthMiddleware(authService))
	{
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/stats", dashboardHandler.GetDashboardStats)
			dashboard.GET("/weekly-chart", dashboardHandler.GetWeeklyChart)
		}

		chargePoints := api.Group("/charge-points")
		{
			chargePoints.GET("", chargePointHandler.GetChargePoints)
			chargePoints.GET("/:id", chargePointHandler.GetChargePoint)
			chargePoints.PATCH("/:id/status", chargePointHandler.UpdateChargePointStatus)
			chargePoints.POST("/:id/commands", chargePointHandler.SendRemoteCommand)
		}

		transactions := api.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
		}

		users := api.Group("/users")
		users.Use(RoleMiddleware("admin"))
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		idTags := api.Group("/id-tags")
		idTags.Use(RoleMiddleware("admin"))
		{
			idTags.GET("", idTagHandler.GetIDTags)
			idTags.GET("/user/:userId", idTagHandler.GetIDTagsByUser)
			idTags.GET("/:id", idTagHandler.GetIDTag)
			idTags.POST("", idTagHandler.CreateIDTag)
			idTags.PUT("/:id", idTagHandler.UpdateIDTag)
			idTags.DELETE("/:id", idTagHandler.DeleteIDTag)
		}
	}
}
