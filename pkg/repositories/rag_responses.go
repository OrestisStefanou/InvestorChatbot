package repositories

import (
	"context"
	"encoding/json"
	"investbot/pkg/services"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ragResponseDocument struct {
	RagTopic     services.Topic
	Conversation []services.Message
	Response     string
	CreatedAt    time.Time
}

type RagResponsesBadgerRepo struct {
	db *badger.DB
}

func NewRagResponsesBadgerRepo(db *badger.DB) (*RagResponsesBadgerRepo, error) {
	return &RagResponsesBadgerRepo{db: db}, nil
}

func (r *RagResponsesBadgerRepo) StoreRagResponse(
	ragTopic services.Topic,
	conversation []services.Message,
	response string,
) error {
	document := ragResponseDocument{
		RagTopic:     ragTopic,
		Conversation: conversation,
		Response:     response,
		CreatedAt:    time.Now(),
	}

	err := r.db.Update(func(txn *badger.Txn) error {
		documentBytes, err := json.Marshal(document)
		if err != nil {
			return err
		}

		return txn.Set([]byte(time.Now().String()), documentBytes)
	})

	return err
}

type RagResponsesMongoRepo struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

func NewRagResponsesMongoRepo(client *mongo.Client, dbName, collectionName string) (*RagResponsesMongoRepo, error) {
	return &RagResponsesMongoRepo{
		client:         client,
		dbName:         dbName,
		collectionName: collectionName,
	}, nil
}

func (r *RagResponsesMongoRepo) StoreRagResponse(
	ragTopic services.Topic,
	conversation []services.Message,
	response string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := ragResponseDocument{
		RagTopic:     ragTopic,
		Conversation: conversation,
		Response:     response,
		CreatedAt:    time.Now(),
	}

	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, document)
	return err
}
