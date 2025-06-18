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
	// 1. Configurar ambiente de teste seguro
	testDBName := "test_auction_auto_close_" + time.Now().Format("20060102150405")
	os.Setenv("AUCTION_DURATION", "2s")
	defer os.Unsetenv("AUCTION_DURATION")

	// 2. Configurar conexão com autenticação
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use suas credenciais reais aqui!
	credential := options.Credential{
		Username: "admin",
		Password: "admin",
	}

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(credential))

	assert.NoError(t, err)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			t.Logf("Disconnect error: %v", err)
		}
	}()

	database := client.Database(testDBName)
	defer database.Drop(ctx) // Garante limpeza após o teste

	repo := auction.NewAuctionRepository(database)

	// 3. Criar leilão de teste
	auctionEntity, internalErr := auction_entity.CreateAuction(
		"Product Test",
		"Category",
		"Description",
		auction_entity.New,
	)
	assert.Nil(t, internalErr)

	// 4. Inserir leilão
	internalErr = repo.CreateAuction(ctx, auctionEntity)
	assert.Nil(t, internalErr)

	// 5. Verificar estado inicial
	var auctionDB auction.AuctionEntityMongo
	err = repo.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).Decode(&auctionDB)
	assert.NoError(t, err)
	assert.Equal(t, auction_entity.Active, auctionDB.Status)

	// 6. Aguardar expiração com margem de segurança
	logger.Info("Waiting for auction to expire...")
	time.Sleep(3 * time.Second) // AUCTION_DURATION (2s) + 1s buffer

	// 7. Executar fechamento
	repo.CloseExpiredAuctions(ctx)

	// 8. Verificar com retry otimizado
	const maxAttempts = 5
	const retryInterval = 500 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = repo.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).Decode(&auctionDB)
		assert.NoError(t, err)

		if auctionDB.Status == auction_entity.Completed {
			break
		}

		if attempt < maxAttempts {
			logger.Info(fmt.Sprintf("Retry %d/%d - Auction not closed yet", attempt, maxAttempts),
				zap.String("status", fmt.Sprintf("%d", auctionDB.Status)),
				zap.Int64("end_time", auctionDB.EndTime))

			time.Sleep(retryInterval)
			repo.CloseExpiredAuctions(ctx) // Executar novamente
		}
	}

	// 9. Verificações finais
	assert.Equal(t, auction_entity.Completed, auctionDB.Status,
		"Auction should be closed with Completed status")

	// Limpar banco de dados de teste
	database.Collection("auctions").Drop(ctx)
}
