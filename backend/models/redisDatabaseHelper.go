package models

import (
	"context"
	"fmt"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/go-redis/redis/v8"
)

type RedisDatabaseConnection struct {
	DatabaseProvider
	redisClient *redis.Client
	logger      logger.ILogger
}

func test() IDatabase {
	return &RedisDatabaseConnection{}
}

func (r *RedisDatabaseConnection) Init(logger logger.ILogger) error {
	r.logger = logger
	return nil
}

func (r *RedisDatabaseConnection) ConnectToDatabase() error {

	ctx := context.Background()

	r.redisClient = redis.NewClient(&redis.Options{
		Addr:        r.Uri,
		Password:    "",
		DB:          0,
		PoolSize:    10,
		DialTimeout: 5 * time.Second,
	})

	statusCmd := r.redisClient.Ping(ctx)
	if statusCmd.Err() != nil {
		r.IsDBConnected = false
		r.logger.InfoLog(fmt.Sprintf("Could not connect to Redis: %v", statusCmd.Err()))
	}

	r.IsDBConnected = true
	r.logger.InfoLog("Successfully connected to Redis!")
	r.logger.InfoLog(fmt.Sprintf("Ping response: %v", statusCmd.Val()))
	return nil
}

func (r *RedisDatabaseConnection) IsConnected() bool {
	return r.IsDBConnected
}

func (r *RedisDatabaseConnection) DisconnectFromDatabase() error {
	if err := r.redisClient.Close(); err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error closing Redis client: %v", err))
		return err
	}

	fmt.Println("Redis client closed.")
	return nil
}

func (r *RedisDatabaseConnection) AddBillIfNotExists(split *Split) (string, error) {
	return "", nil
}

func (r *RedisDatabaseConnection) DeleteBillIfExists(splitId string) error {
	return nil
}

// admin actions
func (r *RedisDatabaseConnection) UpdateBillMetaData(splitId string, newSplit Split) error {
	return nil
}

func (r *RedisDatabaseConnection) AddBillItem(splitId string, item Item) error {
	return nil
}

func (r *RedisDatabaseConnection) DeleteBillItem(splitId string, itemId string) error {
	return nil
}

// admin + other user actions
func (r *RedisDatabaseConnection) AddItemTaker(splitId string, itemId string, takerId string) error {
	return nil
}

func (r *RedisDatabaseConnection) DeleteItemTaker(splitId string, itemId string, takerId string) error {
	return nil
}

/*
{
  "billId": "unique-bill-link-id",
  "items": [
    {
      "id": "item1",
      "name": "Pizza",
      "price": 25.00,
      "takers": ["userA", "userB"]
    },
    {
      "id": "item2",
      "name": "Coke",
      "price": 3.00,
      "takers": ["userC"]
    }
  ],
  "participants": {
    "userA": "Alice",
    "userB": "Bob",
    "userC": "Charlie",
  },
  "status": "active",
  "location": "place name",
  "date": "date",
  "createdAt": "...",
  "lastUpdated": "..."
}

*/
