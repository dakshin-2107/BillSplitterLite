package managers

import (
	"encoding/json"
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
)

func (s *SessionManager) ExecuteAction(billID string, userID string, actionData []byte) (common.IAction, error) {
	var action Action
	err := json.Unmarshal(actionData, &action)
	if err == nil {

		switch action.ActionType {

		case common.HELLO_THERE:
			s.logger.DebugLog("Hello There")

		case common.PING:
			s.logger.DebugLog("Ping")

		case common.BYE_BYE:
			s.logger.DebugLog("Bye Bye")
			return nil, fmt.Errorf("bye bye")

		// item cases
		case common.ADD_NEW_ITEM:
			newItem := common.Item{
				Id:     action.ItemId,
				Name:   action.ItemName,
				Price:  action.Price,
				Takers: make(map[string]int, 0),
			}

			s.logger.DebugLog(fmt.Sprintf("adding item(itemID: %v)", action.ItemId))
			err = s.ModelHelper.AddBillItem(billID, newItem)

		case common.DELETE_ITEM:
			s.logger.DebugLog(fmt.Sprintf("removing item(itemID: %v)", action.ItemId))
			err = s.ModelHelper.DeleteBillItem(billID, action.ItemId)

		case common.EDIT_ITEM:
			s.logger.DebugLog(fmt.Sprintf("editing item(itemID: %v)", action.ItemId))

		// item's taker cases
		case common.ADD_TAKER_FOR_ITEM:
			s.logger.DebugLog(fmt.Sprintf("adding taker(takerId: %v) for item(itemId : %v)", action.TakerId, action.ItemId))
			err = s.ModelHelper.AddItemTaker(billID, action.ItemId, action.TakerId)

		case common.DELETE_TAKER_FOR_ITEM:
			s.logger.DebugLog(fmt.Sprintf("removing taker(takerId: %v) for item(itemId : %v)", action.TakerId, action.ItemId))
			err = s.ModelHelper.DeleteItemTaker(billID, action.ItemId, action.TakerId)

		case common.INCREMENT_TAKER_ID:
			s.logger.DebugLog(fmt.Sprintf("incrementing taker(takerID: %v)", action.TakerId))

		case common.DECREMENT_TAKER_ID:
			s.logger.DebugLog(fmt.Sprintf("decrementing taker(takerID: %v)", action.TakerId))

		// taker cases
		case common.ADD_NEW_TAKER:
			s.logger.DebugLog(fmt.Sprintf("adding taker(takerID: %v) with name : %v", action.TakerId, action.ItemName))
			err = s.ModelHelper.AddTakerID(billID, action.TakerId, action.ItemName)

		case common.DELETE_TAKER:
			s.logger.DebugLog(fmt.Sprintf("deleting taker(takerID: %v)", action.TakerId))
			err = s.ModelHelper.DeleteTakerID(billID, action.TakerId)

		case common.EDIT_TAKER:
			s.logger.DebugLog(fmt.Sprintf("editing taker(takerID: %v)", action.TakerId))
		}
	}

	return &action, err
}
