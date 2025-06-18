package bid

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (br *BidRepository) FindBidByAuctionId(
	ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {

	filter := bson.M{"auction_id": auctionId}
	cursor, err := br.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding bids", err)
		return []bid_entity.Bid{}, nil // Retornar array vazio
	}
	defer cursor.Close(ctx)

	var bids []bid_entity.Bid
	for cursor.Next(ctx) {
		var bidMongo BidEntityMongo
		if err := cursor.Decode(&bidMongo); err != nil {
			logger.Error("Error decoding bid", err)
			continue
		}

		bids = append(bids, bid_entity.Bid{
			Id:        bidMongo.Id,
			UserId:    bidMongo.UserId,
			AuctionId: bidMongo.AuctionId,
			Amount:    bidMongo.Amount,
			Timestamp: time.Unix(bidMongo.Timestamp, 0),
		})
	}

	return bids, nil
}

func (bd *BidRepository) FindWinningBidByAuctionId(
	ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}

	var bidEntityMongo BidEntityMongo
	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		logger.Error("Error trying to find the auction winner", err)
		return nil, internal_error.NewInternalServerError("Error trying to find the auction winner")
	}

	return &bid_entity.Bid{
		Id:        bidEntityMongo.Id,
		UserId:    bidEntityMongo.UserId,
		AuctionId: bidEntityMongo.AuctionId,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
