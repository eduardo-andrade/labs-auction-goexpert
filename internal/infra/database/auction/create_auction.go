package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
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

func (ar *AuctionRepository) FindAuctionById(
	ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, internal_error.NewBadRequestError("Invalid auction ID format")
	}

	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, internal_error.NewNotFoundError("Auction not found")
		}
		logger.Error(fmt.Sprintf("Error trying to find auction by id = %s", id), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
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

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error inserting auction", err)
		return internal_error.NewInternalServerError("Error inserting auction")
	}

	logger.Info("Auction created",
		zap.String("id", auctionEntity.Id),
		zap.Int64("end_time", endTime),
		zap.String("duration", duration.String()))

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
	ar.CloseExpiredAuctions(ctx)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ar.CloseExpiredAuctions(ctx)
		}
	}
}

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) {
	ar.closeMutex.Lock()
	defer ar.closeMutex.Unlock()

	now := time.Now().Unix()
	logger.Info("Checking for expired auctions", zap.Int64("now", now))

	filter := bson.M{
		"status":   auction_entity.Active,
		"end_time": bson.M{"$lte": now},
	}

	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error closing auctions", err)
		return
	}

	if result.ModifiedCount > 0 {
		logger.Info("Closed auctions",
			zap.Int64("count", result.ModifiedCount),
			zap.Any("filter", filter))
	} else {
		logger.Info("No auctions to close")
	}
}
