package auction_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *AuctionController) FindAuctions(c *gin.Context) {
	statusStr := c.Query("status")
	category := c.Query("category")
	productName := c.Query("productName")

	var status auction_entity.AuctionStatus = auction_entity.Active
	if statusStr != "" {
		statusNumber, errConv := strconv.Atoi(statusStr)
		if errConv != nil {
			restErr := rest_err.NewBadRequestError("Invalid status value")
			c.JSON(restErr.Code, restErr)
			return
		}
		status = auction_entity.AuctionStatus(statusNumber)
	}

	auctions, err := u.findUseCase.FindAuctions(
		context.Background(),
		status,
		category,
		productName,
	)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (u *AuctionController) FindAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := rest_err.NewBadRequestError("Invalid auction ID")
		c.JSON(restErr.Code, restErr)
		return
	}

	auctionData, err := u.findUseCase.FindAuctionById(context.Background(), auctionId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		if restErr.Code == http.StatusNotFound {
			c.JSON(restErr.Code, gin.H{"error": "Auction not found"})
			return
		}

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (u *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := rest_err.NewBadRequestError("Invalid auction ID")
		c.JSON(restErr.Code, restErr)
		return
	}

	bid, err := u.findUseCase.FindWinningBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, bid)
}
