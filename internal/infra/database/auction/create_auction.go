package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	// Set the auction status to Active and calculate the expiration timestamp
	durationStr := os.Getenv("AUCTION_DURATION")
	if durationStr == "" {
		logger.Info("AUCTION_DURATION is not set")
		return internal_error.NewInternalServerError("AUCTION_DURATION is not set")
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.Error("Invalid auction duration", err)
		return internal_error.NewInternalServerError("Invalid auction duration")
	}

	auctionEntity.Timestamp = time.Now().Add(time.Duration(duration) * time.Second)

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auction_entity.Active,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err = ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// Start a goroutine to monitor the auction and close it after the duration
	go ar.monitorAuction(ctx, auctionEntityMongo)

	return nil
}

func (ar *AuctionRepository) monitorAuction(ctx context.Context, auction *AuctionEntityMongo) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Now().Unix() > auction.Timestamp {
				auction.Status = auction_entity.Completed
				filter := bson.M{"_id": auction.Id}
				update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
				_, err := ar.Collection.UpdateOne(ctx, filter, update)
				if err != nil {
					logger.Error("Failed to close auction", err, zap.Error(err))
				} else {
					logger.Info("Auction closed successfully", zap.String("auctionID", auction.Id))
				}
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
