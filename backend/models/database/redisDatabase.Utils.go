package database

import (
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
)

func (r *RedisDatabaseConnection) DoesSplitIdExist(splitId string) (bool, error) {
	// Check if the bill already exists
	exists, err := r.redisClient.Exists(r.ctx, splitId).Result()
	return (exists == 1), err
}

func ParseArray[T any](result interface{}) ([]T, error) {
	var requiredArray []T = nil
	var err error = nil
	if arr, ok := result.([]interface{}); ok && len(arr) > 0 {
		if keys, ok := arr[0].([]interface{}); ok {
			requiredArray = make([]T, 0)
			for _, k := range keys {
				if element, ok := k.(T); ok {
					requiredArray = append(requiredArray, element)
				} else {
					err = fmt.Errorf("could not cast invalid type found")
				}
			}
		} else {
			err = fmt.Errorf("no keys found at path")
		}
	} else {
		err = fmt.Errorf("empty or invalid response from Redis")
	}

	return requiredArray, err
}

func (r *RedisDatabaseConnection) DoesItemIdExist(splitId string, itemId string) (bool, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.OBJKEYS", splitId, "$.items").Result()
	if err != nil {
		r.logger.InfoLog(fmt.Sprintf("Error getting items from Redis: %v", err))
		return false, err
	}

	itemIDs, err := ParseArray[string](result)
	if err != nil {
		r.logger.InfoLog("error parsing the response from redis")
		return false, err
	}

	index := -1
	for i, itemID := range itemIDs {
		if itemID == itemId {
			index = i
			break
		}
	}

	return index != -1, nil
}

func testDB1() common.IDatabase {
	return &RedisDatabaseConnection{}
}

func testDB2() common.ISplitModifier {
	return &RedisDatabaseConnection{}
}
