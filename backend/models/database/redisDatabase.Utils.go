package database

import (
	"encoding/json"
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/go-redis/redis/v8"
)

func (r *RedisDatabaseConnection) DoesSplitIdExist(splitId string) (bool, error) {
	exists, err := r.redisClient.Exists(r.ctx, splitId).Result()
	return (exists == 1), err
}

// For OBJKEYS - parses [][]interface{}, returns => []T
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

// For NUMINCRBY - if operation was successful, returns => updatedCount, true
func (r *RedisDatabaseConnection) CheckResult(result interface{}) (int, bool) {
	var arr []int
	if arrStr, ok := result.(string); ok {
		err := json.Unmarshal([]byte(arrStr), &arr)
		if err != nil {
			return -1, false
		}

		if len(arr) > 0 {
			return arr[0], true
		}

	}

	return -1, false
}

func (r *RedisDatabaseConnection) DoesItemIdExist(splitId string, itemId string) (bool, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.OBJKEYS", splitId, "$.items").Result()
	if err != nil {
		return false, err
	}

	itemIDs, err := ParseArray[string](result)
	if err != nil {
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

func (r *RedisDatabaseConnection) DoesTakerIdExist(splitId string, takerId string) (bool, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.OBJKEYS", splitId, "$.participants").Result()
	if err != nil {
		return false, err
	}

	takerIDs, err := ParseArray[string](result)
	if err != nil {
		return false, err
	}

	index := -1
	for i, takerID := range takerIDs {
		if takerID == takerId {
			index = i
			break
		}
	}

	return index != -1, nil
}

// If a refresh is required
func (r *RedisDatabaseConnection) GetBill(billId string) (*common.Split, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.GET", billId, "$").Result()
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
		return nil, fmt.Errorf("bill not found")
	}
	return &splits[0], nil
}

// Checks whether the takerID is present in takers collection of any items in the split
func (r *RedisDatabaseConnection) IsTakerIdUsed(splitID string, takerID string) bool {
	result, err := r.redisClient.Do(r.ctx, "JSON.OBJKEYS", splitID, "$.items").Result()
	if err == nil {
		if itemKeys, err := ParseArray[string](result); err == nil {
			for _, itemKey := range itemKeys {
				itemParticipantsPath := fmt.Sprintf("$.items.%v.takers.%v", itemKey, takerID)
				_, err = r.redisClient.Do(r.ctx, "JSON.GET", splitID, itemParticipantsPath).Result()
				if err != redis.Nil {
					return true
				}
			}
		}
	}

	return false
}

func testDB1() common.IDatabase {
	return &RedisDatabaseConnection{}
}

func testDB2() common.ISplitModifier {
	return &RedisDatabaseConnection{}
}
