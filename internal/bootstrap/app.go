package bootstrapp

import (
	"rentalin/config"
	"rentalin/infrastructure/database"
	"rentalin/internal/handler"
	"rentalin/internal/integration/xendit"
	"rentalin/internal/job"
	"rentalin/internal/repository"
	"rentalin/internal/service"
	"rentalin/pkg/logger"
	"rentalin/pkg/transaction"
	"rentalin/routes"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func NewApp(cfg *config.Config) (*echo.Echo, *pgxpool.Pool) {

	// ================= DATABASE =================
	db := database.NewPosgres(cfg)
	txRunner := transaction.NewTxRunner(db)

	// ================= XENDIT =================
	xenditClient := xendit.NewClient(cfg.XenditSecretKey)

	// ================= REPOSITORY =================
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	rentalRepo := repository.NewRentalRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// ================= SERVICE =================
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	rentalService := service.NewRentalService(txRunner, rentalRepo, productRepo)
	paymentService := service.NewPaymentService(txRunner, paymentRepo, rentalRepo, productRepo, xenditClient)
	dashboardService := service.NewDashboardService(userRepo, productRepo, rentalRepo, paymentRepo)

	// ================= JOB =================
	job.StartPaymentExpiryJob(paymentRepo)
	job.StartOverdueChecker(rentalRepo)

	// ================= HANDLER =================
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	rentalHandler := handler.NewRentalHandler(rentalService)
	paymentHandler := handler.NewPaymentHandler(cfg, paymentService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	logger.Init()

	// ================= ECHO =================
	e := echo.New()

	// ================= ROUTES =================
	routes.SetupRoutes(e, routes.RouteConfig{
		JWTSecret:        cfg.JWTSecret,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		ProductHandler:   productHandler,
		RentalHandler:    rentalHandler,
		PaymentHandler:   paymentHandler,
		DashboardHandler: dashboardHandler,
	})

	return e, db
}
