package managers

import (
	"encoding/json"

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

// dummy method
func (a *Action) ExecuteAction() error {
	return nil
}

func (s *SessionManager) ExecuteAction(billID string, userID string, actionData []byte) (common.IAction, error) {

	/*
	   type Item struct {
	   	ID     string   `json:"id"`
	   	Name   string   `json:"name"`
	   	Price  float64  `json:"price"`
	   	Takers []string `json:"takers"`
	   }
	*/

	var action Action
	err := json.Unmarshal(actionData, &action)
	if err == nil {

		switch action.ActionType {

		// id, name and price must be present
		case common.ADD_ITEM:
			newItem := common.Item{
				Id:     action.ItemId,
				Name:   action.ItemName,
				Price:  action.Price,
				Takers: make([]string, 0),
			}

			s.logger.DebugLog("adding item")
			err = s.ModelHelper.AddBillItem(billID, newItem)
			return &action, err

		// valid id must be present
		case common.REMOVE_ITEM:
			s.logger.DebugLog("removing item")
			err = s.ModelHelper.DeleteBillItem(billID, action.ItemId)
			return &action, err

		// valid id and taker name must be present
		case common.ADD_TAKER:
			s.logger.DebugLog("adding taker")

		case common.REMOVE_TAKER:
			s.logger.DebugLog("removing taker")
		}

	}

	return nil, err
}
