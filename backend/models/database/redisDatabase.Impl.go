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

func (r *RedisDatabaseConnection) AddBillIfNotExists(split *common.Split) error {
	if exists, err := r.DoesSplitIdExist(split.BillID); err == nil && !exists {
		data, err := json.Marshal(split)
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error marshalling split: %v", err))
			return err
		}

		err = r.redisClient.Do(r.ctx, "JSON.SET", split.BillID, "$", data).Err()
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error setting split in Redis: %v", err))
			return err
		}
	}

	return nil
}

func (r *RedisDatabaseConnection) AddBillItem(splitId string, item common.Item) error {
	exist, err := r.DoesItemIdExist(splitId, item.Id)
	if err == nil && !exist {
		data, err := json.Marshal(item)
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error marshalling split: %v", err))
			return err
		}

		path := fmt.Sprintf("$.items.%v", item.Id)
		err = r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, data).Err()
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error appending item to split in Redis: %v", err))
			return err
		}
	} else {
		if err != nil {
			r.logger.InfoLog("error while getting item keys")
		} else {
			r.logger.InfoLog(fmt.Sprintf("item with id %v already exists", item.Id))
		}
	}

	return err
}

func (r *RedisDatabaseConnection) DeleteBillItem(splitId string, itemId string) error {
	if exist, err := r.DoesItemIdExist(splitId, itemId); err == nil && exist {
		path := fmt.Sprintf("$.items.%v", itemId)
		err = r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Err()
		if err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error removing item from split in Redis: %v", err))
			return err
		}
	}

	return nil
}

// Verify that taker exists in the participants dictionary and then add
func (r *RedisDatabaseConnection) AddItemTaker(splitId string, itemId string, takerId string) error {
	// Get items array
	result, err := r.redisClient.Do(r.ctx, "JSON.GET", splitId, "$.items").Result()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error getting items from Redis: %v", err))
		return err
	}

	var items [][]common.Item
	switch res := result.(type) {
	case string:
		if err := json.Unmarshal([]byte(res), &items); err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error unmarshalling items: %v", err))
			return err
		}
	case []byte:
		if err := json.Unmarshal(res, &items); err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error unmarshalling items: %v", err))
			return err
		}
	default:
		r.logger.InfoLog("Unexpected type from Redis JSON.GET for items")
		return fmt.Errorf("unexpected type from Redis JSON.GET for items")
	}

	if len(items) == 0 {
		return fmt.Errorf("no items found")
	}

	index := -1
	for i, item := range items[0] {
		if item.Id == itemId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("item not found")
	}

	// Append takerId to the takers array of the item at index
	path := fmt.Sprintf("$.items[%d].takers", index)
	err = r.redisClient.Do(r.ctx, "JSON.ARRAPPEND", splitId, path, fmt.Sprintf("\"%s\"", takerId)).Err()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error appending taker to item in Redis: %v", err))
		return err
	}

	return nil
}

// Verify that taker exists in the participants dictionary and then delete
func (r *RedisDatabaseConnection) DeleteItemTaker(splitId string, itemId string, takerId string) error {
	result, err := r.redisClient.Do(r.ctx, "JSON.GET", splitId, "$.items").Result()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error getting items from Redis: %v", err))
		return err
	}

	var items [][]common.Item
	switch res := result.(type) {
	case string:
		if err := json.Unmarshal([]byte(res), &items); err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error unmarshalling items: %v", err))
			return err
		}
	case []byte:
		if err := json.Unmarshal(res, &items); err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error unmarshalling items: %v", err))
			return err
		}
	default:
		r.logger.InfoLog("Unexpected type from Redis JSON.GET for items")
		return fmt.Errorf("unexpected type from Redis JSON.GET for items")
	}

	if len(items) == 0 {
		return fmt.Errorf("no items found")
	}

	index := -1
	for i, item := range items[0] {
		if item.Id == itemId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("item not found")
	}

	takerIdx := -1
	for i, t := range items[0][index].Takers {
		if t == takerId {
			takerIdx = i
			break
		}
	}

	if takerIdx == -1 {
		return fmt.Errorf("taker not found")
	}

	path := fmt.Sprintf("$.items[%d].takers", index)
	err = r.redisClient.Do(r.ctx, "JSON.ARRPOP", splitId, path, takerIdx).Err()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error removing taker from item in Redis: %v", err))
		return err
	}

	return nil
}

// If a refresh is required
func (r *RedisDatabaseConnection) GetBill(billId string) (*common.Split, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.GET", billId, "$").Result()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error getting bill from Redis: %v", err))
		return nil, err
	}

	var splits []common.Split
	switch res := result.(type) {
	case string:
		if err := json.Unmarshal([]byte(res), &splits); err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error unmarshalling split: %v", err))
			return nil, err
		}
	case []byte:
		if err := json.Unmarshal(res, &splits); err != nil {
			r.logger.InfoLog(fmt.Sprintf("Error unmarshalling split: %v", err))
			return nil, err
		}
	default:
		r.logger.InfoLog("Unexpected type from Redis JSON.GET")
		return nil, fmt.Errorf("unexpected type from Redis JSON.GET")
	}

	if len(splits) == 0 {
		return nil, fmt.Errorf("bill not found")
	}
	return &splits[0], nil
}

// maybe wont be required
func (r *RedisDatabaseConnection) DeleteBillIfExists(splitId string) error {
	return nil
}

// admin actions
func (r *RedisDatabaseConnection) UpdateBillMetaData(splitId string, newSplit common.Split) error {
	return nil
}
