package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {

	result := []auction_entity.Auction{}

	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions", err)
		return result, nil // Retornar array vazio
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return result, nil // Retornar array vazio
	}

	for _, value := range auctionsMongo {
		result = append(result, auction_entity.Auction{
			Id:          value.Id,
			ProductName: value.ProductName,
			Category:    value.Category,
			Description: value.Description,
			Condition:   value.Condition,
			Status:      value.Status,
			Timestamp:   time.Unix(value.Timestamp, 0),
		})
	}

	return result, nil
}
