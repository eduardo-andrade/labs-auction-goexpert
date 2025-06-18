package auction_usecase

import (
	"fullcycle-auction_go/internal/entity/auction_entity"
	"time"
)

type (
	AuctionInputDTO struct {
		ProductName string
		Category    string
		Description string
		Condition   auction_entity.ProductCondition
	}

	AuctionOutputDTO struct {
		Id          string
		ProductName string
		Category    string
		Description string
		Condition   auction_entity.ProductCondition
		Status      auction_entity.AuctionStatus
		Timestamp   time.Time
	}

	WinningInfoOutputDTO struct {
		Auction AuctionOutputDTO
		Bid     *BidOutputDTO
	}

	BidOutputDTO struct {
		Id        string
		UserId    string
		AuctionId string
		Amount    float64
		Timestamp time.Time
	}
)
