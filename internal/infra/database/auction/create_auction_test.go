package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestAutomaticAuctionClosure verifica se o leilão é fechado automaticamente
// após o tempo configurado via AUCTION_DURATION_SECONDS.
func TestAutomaticAuctionClosure(t *testing.T) {
	// Configura um tempo de expiração curto para teste
	os.Setenv("AUCTION_DURATION_SECONDS", "1")
	defer os.Unsetenv("AUCTION_DURATION_SECONDS")

	// Conecta ao MongoDB local (ex: via docker-compose)
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skip("⚠️ MongoDB não disponível, teste ignorado:", err)
		return
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_auction_db")
	defer db.Drop(context.Background())

	repo := NewAuctionRepository(db)
	ctx := context.Background()

	auction := &auction_entity.Auction{
		Id:          "test_auction_123",
		ProductName: "Console PlayBox X",
		Category:    "Games",
		Description: "Leilão automático de teste",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	// Cria o leilão no Mongo
	errCreate := repo.CreateAuction(ctx, auction)
	if errCreate != nil {
		t.Fatalf("❌ erro ao criar leilão: %v", errCreate)
	}

	// Aguarda mais que o tempo configurado para o fechamento automático
	time.Sleep(2 * time.Second)

	// Busca o leilão novamente no MongoDB
	var result AuctionEntityMongo
	err = repo.Collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&result)
	if err != nil {
		t.Fatalf("❌ erro ao buscar leilão após fechamento: %v", err)
	}

	// Valida o status
	if result.Status != auction_entity.Completed {
		t.Fatalf("❌ esperado status 'Completed', obtido: %v", result.Status)
	}

	t.Logf("✅ Leilão '%s' foi fechado automaticamente com sucesso.", result.Id)
}
