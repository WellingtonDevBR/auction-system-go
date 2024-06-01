package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB() (*mongo.Database, func()) {
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/?authSource=admin")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	db := client.Database("auction_test")
	cleanup := func() {
		db.Drop(context.TODO())
		client.Disconnect(context.TODO())
	}

	return db, cleanup
}

func decodeAuction(singleResult *mongo.SingleResult, auctionInDB *AuctionEntityMongo) *internal_error.InternalError {
	err := singleResult.Decode(auctionInDB)
	if err != nil {
		return internal_error.NewInternalServerError(err.Error())
	}
	return nil
}

func TestCreateAuctionWithAutoClose(t *testing.T) {
	os.Setenv("AUCTION_DURATION", "2") // 2 seconds for quick testing
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewAuctionRepository(db)
	auction, _ := auction_entity.CreateAuction("Test Product", "Category", "Description", auction_entity.New)

	err := repo.CreateAuction(context.Background(), auction)
	assert.Nil(t, err)

	// Wait for the auction to be automatically closed
	time.Sleep(5 * time.Second) // Aumentado para garantir tempo suficiente

	filter := bson.M{"_id": auction.Id}
	var auctionInDB AuctionEntityMongo

	singleResult := repo.Collection.FindOne(context.Background(), filter)
	if singleResult.Err() != nil {
		assert.Fail(t, "Failed to find auction", singleResult.Err().Error())
	}

	internalErr := decodeAuction(singleResult, &auctionInDB)
	assert.Nil(t, internalErr)
	assert.Equal(t, auction_entity.Completed, auctionInDB.Status)
}
