package repositories

import (
	"context"
	"encoding/json"
	"investbot/pkg/services"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type topicAndTagsDocument struct {
	Topic     services.Topic
	Tags      services.Tags
	Question  string
	SessionID string
	UserID    string
	CreatedAt time.Time
}

type TopicAndTagsBagderRepo struct {
	db *badger.DB
}

func NewTopicAndTagsBagderRepo(db *badger.DB) (*TopicAndTagsBagderRepo, error) {
	return &TopicAndTagsBagderRepo{db: db}, nil
}

func (r *TopicAndTagsBagderRepo) StoreTopicAndTags(
	topic services.Topic,
	tags services.Tags,
	question string,
	sessionID string,
	userID string,
) error {
	document := topicAndTagsDocument{
		Topic:     topic,
		Tags:      tags,
		Question:  question,
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	err := r.db.Update(func(txn *badger.Txn) error {
		documentBytes, err := json.Marshal(document)
		if err != nil {
			return err
		}

		return txn.Set([]byte(document.Question), documentBytes)
	})

	return err
}

type TopicAndTagsMongoRepo struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

func NewTopicAndTagsMongoRepo(client *mongo.Client, dbName, collectionName string) (*TopicAndTagsMongoRepo, error) {
	return &TopicAndTagsMongoRepo{
		client:         client,
		dbName:         dbName,
		collectionName: collectionName,
	}, nil
}

func (r *TopicAndTagsMongoRepo) StoreTopicAndTags(
	topic services.Topic,
	tags services.Tags,
	question string,
	sessionID string,
	userID string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := topicAndTagsDocument{
		Topic:     topic,
		Tags:      tags,
		Question:  question,
		SessionID: sessionID,
		UserID:    userID,
	}

	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, document)
	return err
}
