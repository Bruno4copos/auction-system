package auction

import (
	"context"
	"os"
	"strconv"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuctionEntityMongo representa o documento salvo no MongoDB
type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

// AuctionRepository fornece acesso à coleção "auctions"
type AuctionRepository struct {
	Collection *mongo.Collection
}

// NewAuctionRepository cria novo repositório
func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

// CreateAuction insere o leilão no banco e agenda fechamento automático.
func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// Agenda fechamento automático em goroutine separada (context.Background para independência)
	go ar.scheduleAuctionClosure(context.Background(), auctionEntityMongo.Id)

	logger.Info("Auction created and auto-close scheduled successfully. id=" + auctionEntityMongo.Id)
	return nil
}

// scheduleAuctionClosure aguarda a duração e tenta fechar o leilão se ainda estiver ativo.
func (ar *AuctionRepository) scheduleAuctionClosure(ctx context.Context, auctionID string) {
	duration := getAuctionDuration()

	logger.Info("Scheduling automatic close for auction " + auctionID + " after " + duration.String())

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		logger.Info("Context canceled before auction auto-close for ID: " + auctionID)
		return
	case <-timer.C:
		ar.closeExpiredAuction(context.Background(), auctionID)
	}
}

// closeExpiredAuction faz o update para marcar o leilão como Completed somente se estiver Active.
func (ar *AuctionRepository) closeExpiredAuction(ctx context.Context, auctionID string) {
	filter := bson.M{"_id": auctionID, "status": auction_entity.Active}
	update := bson.M{"$set": bson.M{
		"status":    auction_entity.Completed,
		"timestamp": time.Now().Unix(),
	}}

	result, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to close expired auction", err)
		return
	}

	if result.ModifiedCount > 0 {
		logger.Info("Auction " + auctionID + " closed automatically.")
	} else {
		logger.Info("Auction " + auctionID + " already closed or not found.")
	}
}

// getAuctionDuration lê AUCTION_DURATION_SECONDS (segundos). Default = 60s.
func getAuctionDuration() time.Duration {
	env := os.Getenv("AUCTION_DURATION_SECONDS")
	if env == "" {
		env = "60"
	}

	seconds, err := strconv.Atoi(env)
	if err != nil || seconds <= 0 {
		seconds = 60
	}

	return time.Duration(seconds) * time.Second
}
