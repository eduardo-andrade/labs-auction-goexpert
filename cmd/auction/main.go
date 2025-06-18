package main

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/infra/database/bid"
	"fullcycle-auction_go/internal/infra/database/user"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()

	log.Println("Starting application...")
	log.Println("MONGODB_URL:", os.Getenv("MONGODB_URL"))
	log.Println("AUCTION_DURATION:", os.Getenv("AUCTION_DURATION"))

	if err := godotenv.Load("/app/.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning: No .env file found, using environment variables")
		}
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err.Error())
		return
	}

	log.Println("Successfully connected to MongoDB")

	router := gin.Default()

	userController, bidController, auctionsController := initDependencies(databaseConnection)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	auctionCreateUseCase := auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository)
	auctionFindUseCase := auction_usecase.NewAuctionFindUseCase(auctionRepository, bidRepository)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))

	auctionController = auction_controller.NewAuctionController(
		auctionCreateUseCase,
		auctionFindUseCase,
	)

	bidController = bid_controller.NewBidController(
		bid_usecase.NewBidUseCase(bidRepository))

	return
}
