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

var errDefault error = fmt.Errorf("something went wrong") // since returning nil means the operation was successful

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

	r.logger.InfoLog("Redis client closed.")
	return errDefault
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
	if err == nil {
		if !exist {
			data, err := json.Marshal(item)
			if err != nil {
				r.logger.InfoLog(fmt.Sprintf("Error marshalling split: %v", err))
				return err
			}

			path := fmt.Sprintf("$.items.%v", item.Id)
			err = r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, data).Err()
			if err != nil {
				r.logger.InfoLog(fmt.Sprintf("error appending item to split in Redis: %v", err))
				return err
			}
		} else {
			err = fmt.Errorf("item already exists")
		}
	}

	return err
}

func (r *RedisDatabaseConnection) DeleteBillItem(splitId string, itemId string) error {
	path := fmt.Sprintf("$.items.%v", itemId)
	result, err := r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Result()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error removing item from split in Redis: %v", err))
		return err
	}

	if successInt, ok := result.(int64); ok {
		if successInt == 0 {
			return fmt.Errorf("item(itemId : %v) was not present ", itemId)
		} else {
			return nil
		}
	}

	return errDefault
}

func (r *RedisDatabaseConnection) AddItemTaker(splitId string, itemId string, takerId string) error {
	if exists, err := r.DoesTakerIdExist(splitId, takerId); exists && err == nil {
		path := fmt.Sprintf("$.items.%v.takers.%v", itemId, takerId)
		result, incrErr := r.redisClient.Do(r.ctx, "JSON.NUMINCRBY", splitId, path, 1).Result()
		if incrErr == nil {
			if updatedCount, incrSuccess := r.CheckResult(result); incrSuccess {
				r.logger.DebugLog(fmt.Sprintf("taker(ID : %v, count : %v) incremented for item(ID : %v)", takerId, updatedCount, itemId))
				return nil
			} else {
				_, setErr := r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, 1).Result()
				if setErr != redis.Nil {
					r.logger.InfoLog(fmt.Sprintf("Error adding item taker to split in Redis: %v", setErr))
					return setErr
				}

				return nil
			}
		}
	}

	return errDefault
}

func (r *RedisDatabaseConnection) DeleteItemTaker(splitId string, itemId string, takerId string) error {
	if exists, err := r.DoesTakerIdExist(splitId, takerId); exists && err == nil {
		path := fmt.Sprintf("$.items.%v.takers.%v", itemId, takerId)
		result, incrErr := r.redisClient.Do(r.ctx, "JSON.NUMINCRBY", splitId, path, -1).Result()
		if incrErr == nil {
			if updatedCount, ok := r.CheckResult(result); ok {
				if updatedCount <= 0 {
					// Delete the taker from the item if share count reaches 0 or less
					err := r.redisClient.Do(r.ctx, "JSON.DEL", splitId, path).Err()
					if err != nil {
						r.logger.InfoLog(fmt.Sprintf("Error deleting zero-share taker from split in Redis: %v", err))
						return err
					}
					r.logger.DebugLog(fmt.Sprintf("taker(ID : %v) removed from item(ID : %v) as count reached %v", takerId, itemId, updatedCount))
				} else {
					r.logger.DebugLog(fmt.Sprintf("taker(ID : %v, count : %v) decrement for item(ID : %v)", takerId, updatedCount, itemId))
				}
				return nil
			}
		} else {
			return incrErr
		}
	}

	return errDefault
}

func (r *RedisDatabaseConnection) AddTakerID(splitId string, takerId string, takerName string) error {
	path := fmt.Sprintf("$.participants.%v", takerId)
	err := r.redisClient.Do(r.ctx, "JSON.SET", splitId, path, fmt.Sprintf("\"%s\"", takerName)).Err()
	if err != redis.Nil {
		r.logger.InfoLog(fmt.Sprintf("Error adding taker to split in Redis: %v", err))
		return err
	}

	return errDefault
}

func (r *RedisDatabaseConnection) DeleteTakerID(splitId string, takerId string) error {
	// check if the taker ID is being used before deleting the ID
	if r.IsTakerIdUsed(splitId, takerId) {
		return fmt.Errorf("cannot delete since taker is part of the split")
	} else {
		participantPath := fmt.Sprintf("$.participants.%v", takerId)
		result, err := r.redisClient.Do(r.ctx, "JSON.DEL", splitId, participantPath).Result()
		if resultInt, ok := result.(int64); ok && err == nil && resultInt == 1 {
			return nil
		}
	}

	return errDefault
}

// maybe wont be required
func (r *RedisDatabaseConnection) DeleteBillIfExists(splitId string) error {
	return nil
}

// admin actions
func (r *RedisDatabaseConnection) UpdateBillMetaData(splitId string, newSplit common.Split) error {
	return nil
}
