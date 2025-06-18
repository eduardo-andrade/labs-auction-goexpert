package bid_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *BidController) FindBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := rest_err.NewBadRequestError("Invalid auction ID")
		c.JSON(restErr.Code, restErr)
		return
	}

	bidOutputList, err := u.bidUseCase.FindBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	if bidOutputList == nil {
		c.JSON(http.StatusOK, []interface{}{})
	} else {
		c.JSON(http.StatusOK, bidOutputList)
	}
}
