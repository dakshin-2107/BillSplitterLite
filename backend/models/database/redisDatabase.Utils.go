package database

import (
	"encoding/json"
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
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

func (r *RedisDatabaseConnection) DoesTakerIdExist(splitId string, takerId string) (bool, error) {
	result, err := r.redisClient.Do(r.ctx, "JSON.OBJKEYS", splitId, "$.participants").Result()
	if err != nil {
		return false, err
	}

	takerIDs, err := ParseArray[string](result)
	if err != nil {
		return false, err
	}

	for _, id := range takerIDs {
		if id == takerId {
			return true, nil
		}
	}

	return false, nil
}

// Checks whether the takerID is present in takers collection of any items in any bill of the split
func (r *RedisDatabaseConnection) IsTakerIdUsed(splitID string, takerID string) bool {
	split, err := r.GetSplit(splitID)
	if err != nil {
		return false
	}

	for _, bill := range split.Bills {
		for _, item := range bill.Items {
			if _, exists := item.Takers[takerID]; exists {
				return true
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
