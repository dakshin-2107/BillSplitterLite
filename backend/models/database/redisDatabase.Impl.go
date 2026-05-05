package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"

	"github.com/go-redis/redis/v8"
)

const ACTIONS_QUEUE_KEY = "ActionsQueue"

var errDefault error = fmt.Errorf("something went wrong")

type RedisDatabaseConnection struct {
	common.DatabaseProvider
	ctx         context.Context
	redisClient *redis.Client
	logger      logger.ILogger
}

func (r *RedisDatabaseConnection) Init(logger logger.ILogger) error {
	r.logger = logger
	r.Uri = os.Getenv("REDIS_URI")
	r.ConnectToDatabase()
	return nil
}

func (r *RedisDatabaseConnection) ConnectToDatabase() error {
	r.ctx = context.Background()
	r.redisClient = redis.NewClient(&redis.Options{
		Addr:        r.Uri,
		Password:    "",
		DB:          0,
		PoolSize:    10,
		DialTimeout: 5 * time.Second,
	})

	statusCmd := r.redisClient.Ping(r.ctx)
	if err := statusCmd.Err(); err != nil {
		r.IsDBConnected = false
		r.logger.InfoLog(fmt.Sprintf("Could not connect to Redis: %v", err))
		return err
	}

	r.IsDBConnected = true
	r.logger.InfoLog("Successfully connected to Redis!")
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
	r.IsDBConnected = false
	r.logger.InfoLog("Redis client closed.")
	return nil
}

// ISplitModifier implementation

func (r *RedisDatabaseConnection) GetSplit(splitId string) (*common.Split, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.GET", splitId, "$").Result()
	if err != nil {
		return nil, err
	}

	var splits []common.Split
	switch res := result.(type) {
	case string:
		if err := json.Unmarshal([]byte(res), &splits); err != nil {
			return nil, err
		}
	case []byte:
		if err := json.Unmarshal(res, &splits); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected type from Redis JSON.GET")
	}

	if len(splits) == 0 {
		return nil, fmt.Errorf("split not found")
	}
	return &splits[0], nil
}

func (r *RedisDatabaseConnection) AddSplitIfNotExists(split *common.Split) error {
	if exists, err := r.DoesSplitIdExist(split.SplitID); err == nil && !exists {
		data, err := json.Marshal(split)
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error marshalling split: %v", err))
			return err
		}

		err = r.redisClient.Do(r.ctx, "JSON.SET", split.SplitID, "$", data).Err()
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error setting split in Redis: %v", err))
			return err
		}
	}
	return nil
}

func (r *RedisDatabaseConnection) DeleteSplitIfExists(splitId string) error {
	err := r.redisClient.Del(r.ctx, splitId).Err()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error deleting split from Redis: %v", err))
		return err
	}
	return nil
}

func (r *RedisDatabaseConnection) GetNewBillId(splitId string) int {
	result, err := r.redisClient.Do(r.ctx, "JSON.NUMINCRBY", splitId, "$.billIdCounter", 1).Result()
	if err == nil {
		if val, ok := r.CheckResult(result); ok {
			return val
		}
	}
	r.logger.InfoLog(fmt.Sprintf("Error getting new bill ID: %v", err))
	return 0
}

func (r *RedisDatabaseConnection) AddBillToSplit(splitId string, bill common.Bill) error {
	data, err := json.Marshal(bill)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("$.bills.%d", bill.BillID)
	return r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, data).Err()
}

func (r *RedisDatabaseConnection) DeleteBillFromSplit(splitId string, billId int) error {
	path := fmt.Sprintf("$.bills.%d", billId)
	result, err := r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Result()
	if err != nil {
		return err
	}
	if successInt, ok := result.(int64); ok && successInt > 0 {
		return nil
	}
	return fmt.Errorf("bill not found")
}

func (r *RedisDatabaseConnection) AddNewTakerToSplit(splitId string, takerId string, takerName string) error {
	path := fmt.Sprintf("$.participants.%v", takerId)
	return r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, fmt.Sprintf("\"%s\"", takerName)).Err()
}

func (r *RedisDatabaseConnection) DeleteTakerFromSplit(splitId string, takerId string) error {
	if r.IsTakerIdUsed(splitId, takerId) {
		return fmt.Errorf("cannot delete since taker is part of a bill")
	}
	path := fmt.Sprintf("$.participants.%v", takerId)
	result, err := r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Result()
	if successInt, ok := result.(int64); ok && err == nil && successInt == 1 {
		return nil
	}
	return fmt.Errorf("participant not found or could not be deleted")
}

func (r *RedisDatabaseConnection) GetNewItemId(splitId string, billId int) int {
	path := fmt.Sprintf("$.bills.%d.itemIdCounter", billId)
	result, err := r.redisClient.Do(r.ctx, "JSON.NUMINCRBY", splitId, path, 1).Result()
	if err == nil {
		if val, ok := r.CheckResult(result); ok {
			return val
		}
	}
	r.logger.InfoLog(fmt.Sprintf("Error getting new item ID: %v", err))
	return 0
}

func (r *RedisDatabaseConnection) AddItemToBill(splitId string, billId int, item common.Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("$.bills.%d.items.%d", billId, item.Id)
	return r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, data).Err()
}

func (r *RedisDatabaseConnection) DeleteItemFromBill(splitId string, billId int, itemId int) error {
	path := fmt.Sprintf("$.bills.%d.items.%d", billId, itemId)
	result, err := r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Result()
	if err != nil {
		return err
	}
	if successInt, ok := result.(int64); ok && successInt > 0 {
		return nil
	}
	return fmt.Errorf("item not found")
}

func (r *RedisDatabaseConnection) EditItemInBill(splitId string, billId int, itemId int, name string, price float32) error {
	basePath := fmt.Sprintf("$.bills.%d.items.%d", billId, itemId)

	nameJSON, err := json.Marshal(name)
	if err != nil {
		return err
	}
	if err := r.redisClient.Do(r.ctx, "JSON.SET", splitId, basePath+".name", string(nameJSON)).Err(); err != nil {
		return err
	}

	rounded := fmt.Sprintf("%.2f", price)
	return r.redisClient.Do(r.ctx, "JSON.SET", splitId, basePath+".price", rounded).Err()
}

func (r *RedisDatabaseConnection) AddNewTakerForItem(splitId string, billId int, itemId int, takerId string) error {
	if exists, err := r.DoesTakerIdExist(splitId, takerId); exists && err == nil {
		path := fmt.Sprintf("$.bills.%d.items.%d.takers.%s", billId, itemId, takerId)
		result, incrErr := r.redisClient.Do(r.ctx, "JSON.NUMINCRBY", splitId, path, 1).Result()
		if incrErr == nil {
			if _, ok := r.CheckResult(result); ok {
				return nil
			} else {
				// Initialize taker with 1 share if it doesn't exist
				return r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, 1).Err()
			}
		}
		return incrErr
	}
	return fmt.Errorf("taker does not exist in split participants")
}

func (r *RedisDatabaseConnection) DeleteTakerForItem(splitId string, billId int, itemId int, takerId string) error {
	path := fmt.Sprintf("$.bills.%d.items.%d.takers.%s", billId, itemId, takerId)
	result, incrErr := r.redisClient.Do(r.ctx, "JSON.NUMINCRBY", splitId, path, -1).Result()
	if incrErr == nil {
		if updatedCount, ok := r.CheckResult(result); ok {
			if updatedCount <= 0 {
				return r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Err()
			}
			return nil
		}
	}
	return incrErr
}

func (r *RedisDatabaseConnection) UpdateBillInformation(splitId string, billId int, newTotal float32, newLocation string, newDate string) error {
	billPath := fmt.Sprintf("$.bills.%d", billId)
	// Update total
	err := r.redisClient.Do(r.ctx, "JSON.SET", splitId, billPath+".total", newTotal).Err()
	if err != nil {
		return err
	}
	// Update location
	err = r.redisClient.Do(r.ctx, "JSON.SET", splitId, billPath+".location", fmt.Sprintf("\"%s\"", newLocation)).Err()
	if err != nil {
		return err
	}
	// Update date
	return r.redisClient.Do(r.ctx, "JSON.SET", splitId, billPath+".date", fmt.Sprintf("\"%s\"", newDate)).Err()
}

func (r *RedisDatabaseConnection) GetBill(splitId string, billId int) (*common.Bill, error) {
	path := fmt.Sprintf("$.bills.%d", billId)
	result, err := r.redisClient.Do(r.ctx, "JSON.GET", splitId, path).Result()
	if err != nil {
		return nil, err
	}

	var bills []common.Bill
	switch res := result.(type) {
	case string:
		err = json.Unmarshal([]byte(res), &bills)
	case []byte:
		err = json.Unmarshal(res, &bills)
	default:
		return nil, fmt.Errorf("unexpected type from Redis JSON.GET")
	}
	if err != nil || len(bills) == 0 {
		return nil, fmt.Errorf("bill not found")
	}
	return &bills[0], nil
}

func (r *RedisDatabaseConnection) AddAllTakersToItem(splitId string, billId int, itemId int) error {
	split, err := r.GetSplit(splitId)
	if err != nil {
		return err
	}

	for takerId := range split.Participants {
		path := fmt.Sprintf("$.bills.%d.items.%d.takers.%s", billId, itemId, takerId)
		// Set taker with 1 share if they don't exist
		err = r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, 1).Err()
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error adding taker %s to item: %v", takerId, err))
		}
	}
	return nil
}

