package auction_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/api/web/validation"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	createUseCase auction_usecase.AuctionUseCaseInterface
	findUseCase   auction_usecase.AuctionFindUseCaseInterface
}

func NewAuctionController(
	createUseCase auction_usecase.AuctionUseCaseInterface,
	findUseCase auction_usecase.AuctionFindUseCaseInterface) *AuctionController {

	return &AuctionController{
		createUseCase: createUseCase,
		findUseCase:   findUseCase,
	}
}

type CreateAuctionRequest struct {
	ProductName string `json:"product_name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description" binding:"required"`
	Condition   string `json:"condition" binding:"required"`
}

func (u *AuctionController) CreateAuction(c *gin.Context) {
	var request CreateAuctionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	// Converter string para ProductCondition
	var condition auction_entity.ProductCondition
	switch strings.ToLower(request.Condition) {
	case "new":
		condition = auction_entity.New
	case "used":
		condition = auction_entity.Used
	case "refurbished":
		condition = auction_entity.Refurbished
	default:
		restErr := rest_err.NewBadRequestError("Invalid condition value. Must be 'new', 'used' or 'refurbished'")
		c.JSON(restErr.Code, restErr)
		return
	}

	auctionInputDTO := auction_usecase.AuctionInputDTO{
		ProductName: request.ProductName,
		Category:    request.Category,
		Description: request.Description,
		Condition:   condition,
	}

	auction, err := u.createUseCase.CreateAuction(context.Background(), auctionInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, auction)
}
