package api
import (
	"database/sql"

	"github.com/KhoirulAziz99/final_project_e_wallet/api/handler"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
	"github.com/KhoirulAziz99/final_project_e_wallet/pkg"
	"github.com/gin-gonic/gin"
)
func SetUpRouter(db *sql.DB) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	userService := app.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userService)
	topupRepo := repository.NewTopupRepository(db)
	topupService := app.NewTopupUsecase(topupRepo)
	topupHandler := handler.NewTopupHandler(topupService)
	walletRepo := repository.NewWalletRepository(db)
	walletService := app.NewWalletUsecase(walletRepo, topupRepo)
	walletHandler := handler.NewWalletHandler(walletService)
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := app.NewPaymentUsecase(paymentRepo, walletRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	transferRepo := repository.NewTransferRepository(db, walletRepo)
	transferService := app.NewTransferUsecase(transferRepo, walletRepo)
	transferHandler := handler.NewTransferHandler(transferService)
	withdrawalRepo := repository.NewWithdrawRepository(db)
	withdrawalService := app.NewWithdrawUsecase(withdrawalRepo, walletRepo)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalService)


	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/history-withdrawal/:walletID", withdrawalHandler.HistoryWithdrawal)
	apiV1.GET("/history-payment/:paymentID", paymentHandler.HistoryPayment)
	apiV1.GET("/history-topup/:topupID", topupHandler.HistoryTopup)
	apiV1.GET("/history-transfer/:sender_id", transferHandler.HistoryTransaction)


	userRouters := apiV1.Group("/users")
	{
		userRouters.POST("/", userHandler.InsertUser)
		userRouters.POST("/login", userHandler.Login)
		 userRouters.Use(pkg.AuthMiddleware())
		userRouters.GET("/profile", handler.ProfileHandler)
		userRouters.PUT("/:id", userHandler.UpdateUser)
		userRouters.DELETE("/:id", userHandler.DeleteUser)
		userRouters.GET("/:id", userHandler.FindOneUser)
		userRouters.GET("/", userHandler.FindAllUsers)
	}

	paymentRouters := apiV1.Group("/payments")
	{	
		paymentRouters.Use(pkg.AuthMiddleware())
		paymentRouters.POST("/", paymentHandler.CreatePayment)
		paymentRouters.GET("/:paymentID", paymentHandler.GetPaymentByID)
		paymentRouters.PUT("/:paymentID", paymentHandler.UpdatePayment)
		paymentRouters.DELETE("/:paymentID", paymentHandler.DeletePayment)
		// paymentRouters.POST("/make-payment", paymentHandler.MakePayment)
	}

	topupRouters := apiV1.Group("/topups")
	{	
		topupRouters.Use(pkg.AuthMiddleware())
		topupRouters.POST("/", topupHandler.CreateTopup)
		topupRouters.GET("/:topupID", topupHandler.GetTopupByID)
		topupRouters.PUT("/:topupID", topupHandler.UpdateTopup)
		topupRouters.DELETE("/:topupID", topupHandler.DeleteTopup)
		topupRouters.GET("/last-amount/:walletID", topupHandler.GetLastTopupAmount)
	}

	transferRouters := apiV1.Group("/transfers")
	{	
		transferRouters.Use(pkg.AuthMiddleware())
		transferRouters.POST("/", transferHandler.CreateTransfer)
		transferRouters.GET("/:transferID", transferHandler.GetTransferByID)
		// transferRouters.PUT("/:transferID", transferHandler.UpdateTransfer)
		transferRouters.DELETE("/:transferID", transferHandler.DeleteTransfer)
		// transferRouters.POST("/make-transfer", transferHandler.MakeTransfer)
	}

	walletRouters := apiV1.Group("/wallets")
	{	
		walletRouters.Use(pkg.AuthMiddleware())
		walletRouters.POST("/", walletHandler.CreateWallet)
		walletRouters.GET("/:walletID", walletHandler.GetWalletByID)
		walletRouters.PUT("/:walletID", walletHandler.UpdateWalletBalance)
		walletRouters.DELETE("/:walletID", walletHandler.DeleteWallet)
	}

	withdrawalRouters := apiV1.Group("/withdrawals")
	{	
		withdrawalRouters.Use(pkg.AuthMiddleware())
		withdrawalRouters.POST("/", withdrawalHandler.CreateWithdrawal)
		withdrawalRouters.GET("/:id", withdrawalHandler.GetWithdrawalByID)
		// withdrawalRouters.PUT("/:id", withdrawalHandler.UpdateWithdrawal)
		withdrawalRouters.DELETE("/:id", withdrawalHandler.DeleteWithdrawal)
		// withdrawalRouters.POST("/:make-withdrawal", withdrawalHandler.MakeWithdrawal)
	}

	return r
}

