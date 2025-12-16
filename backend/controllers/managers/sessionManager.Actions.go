package managers

import (
	"encoding/json"
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
)

func (s *SessionManager) QueueAction(action common.IAction) error {
	return nil
}

func (s *SessionManager) DequeueAction() (common.IAction, error) {

	return nil, nil
}

func (s *SessionManager) PublishAction(action common.IAction) error {

	return nil
}

// dummy method for now
func (a *Action) ExecuteAction() error {
	return nil
}

func (s *SessionManager) ExecuteAction(billID string, userID string, actionData []byte) (common.IAction, error) {
	var action Action
	err := json.Unmarshal(actionData, &action)
	if err == nil {

		switch action.ActionType {
		case common.ADD_ITEM:
			newItem := common.Item{
				Id:     action.ItemId,
				Name:   action.ItemName,
				Price:  action.Price,
				Takers: make(map[string]int, 0),
			}

			s.logger.DebugLog(fmt.Sprintf("adding item(itemID: %v)", action.ItemId))
			err = s.ModelHelper.AddBillItem(billID, newItem)

		case common.REMOVE_ITEM:
			s.logger.DebugLog(fmt.Sprintf("removing item(itemID: %v)", action.ItemId))
			err = s.ModelHelper.DeleteBillItem(billID, action.ItemId)

		case common.ADD_ITEM_TAKER:
			s.logger.DebugLog(fmt.Sprintf("adding taker(takerId: %v) for item(itemId : %v)", action.TakerId, action.ItemId))
			err = s.ModelHelper.AddItemTaker(billID, action.ItemId, action.TakerId)

		case common.REMOVE_ITEM_TAKER:
			s.logger.DebugLog(fmt.Sprintf("removing taker(takerId: %v) for item(itemId : %v)", action.TakerId, action.ItemId))
			err = s.ModelHelper.DeleteItemTaker(billID, action.ItemId, action.TakerId)

		case common.ADD_TAKER_ID:
			s.logger.DebugLog(fmt.Sprintf("adding taker(takerID: %v) with name : %v", action.TakerId, action.ItemName))
			err = s.ModelHelper.AddTakerID(billID, action.TakerId, action.ItemName)

		case common.REMOVE_TAKER_ID:
			s.logger.DebugLog(fmt.Sprintf("removing taker(takerID: %v)", action.TakerId))
			err = s.ModelHelper.DeleteTakerID(billID, action.TakerId)

		case common.BYE_BYE:
			s.logger.DebugLog("Bye Bye")
			return nil, fmt.Errorf("bye bye")
		}
	}

	return &action, err
}
