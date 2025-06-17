package auction_test

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func TestAuctionAutoClose(t *testing.T) {
	os.Setenv("AUCTION_DURATION", "2s")
	defer os.Unsetenv("AUCTION_DURATION")

	// Configurar MongoDB de teste
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	database := client.Database("test_auction_auto_close")
	repo := auction.NewAuctionRepository(database)

	// Criar leilão de teste
	auctionEntity, internalErr := auction_entity.CreateAuction(
		"Product Test",
		"Category",
		"Description",
		auction_entity.New,
	)
	assert.Nil(t, internalErr)

	// Inserir leilão
	internalErr = repo.CreateAuction(ctx, auctionEntity)
	assert.Nil(t, internalErr)

	// Verificar que o leilão está ativo inicialmente
	var auctionDB auction.AuctionEntityMongo
	err = repo.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).Decode(&auctionDB)
	assert.NoError(t, err)
	assert.Equal(t, auction_entity.Active, auctionDB.Status)

	logger.Info("Auction created",
		zap.String("id", auctionEntity.Id),
		zap.Int64("end_time", auctionDB.EndTime))

	// Aguardar tempo de expiração + buffer
	time.Sleep(3 * time.Second)

	// Executar o fechamento de leilões expirados
	repo.CloseExpiredAuctions(ctx)

	// Verificar com retry se o leilão foi fechado
	startTime := time.Now()
	timeout := 5 * time.Second
	closed := false

	for time.Since(startTime) < timeout {
		err = repo.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).Decode(&auctionDB)
		if err != nil {
			logger.Error("Error finding auction", err)
			break
		}

		if auctionDB.Status == auction_entity.Completed {
			closed = true
			break
		}

		logger.Info("Auction not closed yet",
			zap.String("status", fmt.Sprintf("%d", auctionDB.Status)),
			zap.Int64("current_time", time.Now().Unix()),
			zap.Int64("end_time", auctionDB.EndTime))

		// Executar novamente o fechamento
		repo.CloseExpiredAuctions(ctx)
		time.Sleep(500 * time.Millisecond)
	}

	// Verificações finais
	assert.NoError(t, err, "Should find auction without error")
	assert.True(t, closed, "Auction should be closed (status Completed)")
	if closed {
		assert.Equal(t, auction_entity.Completed, auctionDB.Status)
	} else {
		t.Fatalf("Auction not closed after %v: %+v", timeout, auctionDB)
	}

	// Limpar banco de dados de teste
	database.Collection("auctions").Drop(ctx)
}
