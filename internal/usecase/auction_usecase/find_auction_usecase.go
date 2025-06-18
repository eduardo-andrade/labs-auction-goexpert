package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
)

type AuctionFindUseCaseInterface interface {
	FindAuctionById(
		ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)

	FindAuctions(
		ctx context.Context,
		status auction_entity.AuctionStatus,
		category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError)

	FindWinningBidByAuctionId(
		ctx context.Context,
		auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type AuctionFindUseCase struct {
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	bidRepositoryInterface     bid_entity.BidEntityRepository
}

func NewAuctionFindUseCase(
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface,
	bidRepositoryInterface bid_entity.BidEntityRepository) AuctionFindUseCaseInterface {

	return &AuctionFindUseCase{
		auctionRepositoryInterface: auctionRepositoryInterface,
		bidRepositoryInterface:     bidRepositoryInterface,
	}
}

func (au *AuctionFindUseCase) FindAuctionById(
	ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError) {

	auctionEntity, err := au.auctionRepositoryInterface.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AuctionOutputDTO{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp,
	}, nil
}

func (au *AuctionFindUseCase) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError) {

	auctionEntities, err := au.auctionRepositoryInterface.FindAuctions(ctx, status, category, productName)
	if err != nil {
		return nil, err
	}

	var auctionOutputs []AuctionOutputDTO
	for _, value := range auctionEntities {
		auctionOutputs = append(auctionOutputs, AuctionOutputDTO{
			Id:          value.Id,
			ProductName: value.ProductName,
			Category:    value.Category,
			Description: value.Description,
			Condition:   value.Condition,
			Status:      value.Status,
			Timestamp:   value.Timestamp,
		})
	}

	return auctionOutputs, nil
}

func (au *AuctionFindUseCase) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {

	auction, err := au.auctionRepositoryInterface.FindAuctionById(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	if auction.Status != auction_entity.Completed {
		return nil, internal_error.NewBadRequestError("Auction is not completed yet")
	}

	bidWinning, err := au.bidRepositoryInterface.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	if bidWinning == nil {
		return nil, internal_error.NewNotFoundError("No winning bid found for this auction")
	}

	bidOutputDTO := &BidOutputDTO{
		Id:        bidWinning.Id,
		UserId:    bidWinning.UserId,
		AuctionId: bidWinning.AuctionId,
		Amount:    bidWinning.Amount,
		Timestamp: bidWinning.Timestamp,
	}

	auctionOutputDTO := &AuctionOutputDTO{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp,
	}

	return &WinningInfoOutputDTO{
		Auction: *auctionOutputDTO,
		Bid:     bidOutputDTO,
	}, nil
}
