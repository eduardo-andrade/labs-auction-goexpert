package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
	EndTime     int64                           `bson:"end_time"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	closeMutex sync.Mutex
}

// FindAuctionById implements auction_entity.AuctionRepositoryInterface.
func (ar *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	panic("unimplemented")
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}

	go repo.StartAuctionCloser(context.Background())

	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	logger.Info("Creating auction in DB",
		zap.String("id", auctionEntity.Id),
		zap.String("product", auctionEntity.ProductName))

	duration := getAuctionDuration()
	endTime := time.Now().Add(duration).Unix()

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
		EndTime:     endTime,
	}

	result, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error inserting auction",
			err,
			zap.Any("auction", auctionEntityMongo))
		return internal_error.NewInternalServerError("Error inserting auction")
	}
	logger.Info("Auction created successfully",
		zap.Any("inserted_id", result.InsertedID))
	return nil
}

func getAuctionDuration() time.Duration {
	durationStr := os.Getenv("AUCTION_DURATION")
	if durationStr == "" {
		durationStr = "24h"
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("Invalid AUCTION_DURATION format, using default 24h", err)
		return 24 * time.Hour
	}
	return duration
}

func (ar *AuctionRepository) StartAuctionCloser(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ar.CloseExpiredAuctions(ctx)
				time.Sleep(10 * time.Second)
			}
		}
	}()
}

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) {
	ar.closeMutex.Lock()
	defer ar.closeMutex.Unlock()

	now := time.Now().Unix()

	_, err := ar.Collection.UpdateMany(
		ctx,
		bson.M{
			"status":   auction_entity.Active,
			"end_time": bson.M{"$lte": now},
		},
		bson.M{"$set": bson.M{"status": auction_entity.Completed}},
	)

	if err != nil {
		logger.Error("Error closing auctions", err)
		return
	}

	logger.Info("Closed expired auctions")
}
