package server

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/malikkhoiri/csms/internal/application/service"
	"github.com/malikkhoiri/csms/internal/config"
	"github.com/malikkhoiri/csms/internal/domain"
	"github.com/malikkhoiri/csms/internal/handler/http"
	"github.com/malikkhoiri/csms/internal/handler/ws"
	"github.com/malikkhoiri/csms/internal/infrastructure/database"
	"github.com/malikkhoiri/csms/internal/infrastructure/repository"
)

type Server struct {
	router *gin.Engine
	port   string
	config *config.Config

	chargePointService domain.ChargePointService
	transactionService domain.TransactionService
	userService        domain.UserService
	connectorService   domain.ConnectorService
	idTagService       domain.IDTagService
	authService        domain.AuthService
}

func NewServer(cfg *config.Config) (*Server, error) {
	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	postgresDB, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		return nil, err
	}

	chargePointRepo := repository.NewChargePointRepository(postgresDB.DB)
	connectorRepo := repository.NewConnectorRepository(postgresDB.DB)
	transactionRepo := repository.NewTransactionRepository(postgresDB.DB)
	userRepo := repository.NewUserRepository(postgresDB.DB)
	idTagRepo := repository.NewIDTagRepository(postgresDB.DB)

	chargePointService := service.NewChargePointService(chargePointRepo, connectorRepo)
	transactionService := service.NewTransactionService(transactionRepo, chargePointRepo, idTagRepo, cfg.Tariff)
	userService := service.NewUserService(userRepo)
	connectorService := service.NewConnectorService(connectorRepo)
	idTagService := service.NewIDTagService(idTagRepo)
	authService := service.NewAuthService(userRepo, &cfg.JWT)

	return &Server{
		router: router,
		port:   cfg.Server.Port,
		config: cfg,

		chargePointService: chargePointService,
		transactionService: transactionService,
		userService:        userService,
		connectorService:   connectorService,
		idTagService:       idTagService,
		authService:        authService,
	}, nil
}

func (s *Server) SetupRoutes() {
	healthHandler := http.NewHealthHandler()
	ocppHandler := ws.NewOCPPHandler(
		s.chargePointService,
		s.transactionService,
		s.userService,
		s.connectorService,
		s.idTagService,
	)

	s.router.GET("/health", healthHandler.HealthCheck)

	// Setup API routes
	http.SetupRoutes(
		s.router,
		s.chargePointService,
		s.transactionService,
		s.userService,
		s.idTagService,
		s.authService,
	)

	// WebSocket OCPP endpoint
	s.router.GET("/ocpp/*cpID", ocppHandler.HandleWebSocket)
}

func (s *Server) Start() error {
	log.Printf("CSMS server is running on port %s", s.port)
	return s.router.Run(":" + s.port)
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
