package managers

import (
	"encoding/json"
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
)

func (s *SessionManager) ExecuteAction(splitID string, userID string, actionData []byte) error {
	requiresTallyComputation := false
	var action Action
	err := json.Unmarshal(actionData, &action)
	if err == nil {

		// Default BillId to 1 if not provided, for backward compatibility or initial bill
		if action.BillId == 0 {
			action.BillId = 1
		}

		switch action.ActionType {

		case common.HELLO_THERE:
			s.logger.DebugLog("Hello There")

		case common.PING:
			s.logger.DebugLog("Ping")

		case common.BYE_BYE:
			s.logger.DebugLog("Bye Bye")
			err = s.ModelHelper.DeleteSplitIfExists(splitID)
			if err == nil {
				if sess, ok := s.getSession(splitID); ok {
					sess.BroadcastChannel <- common.ActionExecutionSuccessResponse(&action)
					sess.SignalChannel <- action
				}
				return nil
			}

		case common.SYNC_BILL_STATE:
			s.logger.DebugLog("Sync Bill State")
			split, err := s.ModelHelper.GetSplit(splitID)
			if err != nil {
				return err
			}

			requiresTallyComputation = true
			if sess, ok := s.getSession(splitID); ok {
				sess.BroadcastChannel <- common.SplitResponse(*split)
			}

		// item cases
		case common.ADD_NEW_ITEM:
			newItem := common.Item{
				Id:     s.ModelHelper.GetNewItemId(splitID, action.BillId),
				Name:   action.ItemName,
				Price:  action.Price,
				Takers: make(map[string]int),
			}

			s.logger.DebugLog(fmt.Sprintf("adding item(itemID: %v) to bill(billID: %v) in split(splitID: %v)", newItem.Id, action.BillId, splitID))
			err = s.ModelHelper.AddItemToBill(splitID, action.BillId, newItem)
			requiresTallyComputation = true

		case common.DELETE_ITEM:
			s.logger.DebugLog(fmt.Sprintf("removing item(itemID: %v) from bill(billID: %v)", action.ItemId, action.BillId))
			err = s.ModelHelper.DeleteItemFromBill(splitID, action.BillId, action.ItemId)
			requiresTallyComputation = true

		case common.EDIT_ITEM:
			s.logger.DebugLog(fmt.Sprintf("editing item(itemID: %v)", action.ItemId))
			// TODO: Implement EditItem in ModelHelper if needed
			// requiresTallyComputation = true only if price changes

		// item's taker cases
		case common.ADD_TAKER_FOR_ITEM:
			s.logger.DebugLog(fmt.Sprintf("adding taker(takerId: %v) for item(itemId : %v) in bill(billId: %v)", action.TakerId, action.ItemId, action.BillId))
			err = s.ModelHelper.AddNewTakerForItem(splitID, action.BillId, action.ItemId, action.TakerId)
			requiresTallyComputation = true

		case common.DELETE_TAKER_FOR_ITEM:
			s.logger.DebugLog(fmt.Sprintf("removing taker(takerId: %v) for item(itemId : %v) in bill(billId: %v)", action.TakerId, action.ItemId, action.BillId))
			err = s.ModelHelper.DeleteTakerForItem(splitID, action.BillId, action.ItemId, action.TakerId)
			requiresTallyComputation = true

		case common.ADD_ALL_TAKERS_FOR_ITEM:
			s.logger.DebugLog(fmt.Sprintf("adding all participants to item(itemId : %v) in bill(billId: %v)", action.ItemId, action.BillId))
			err = s.ModelHelper.AddAllTakersToItem(splitID, action.BillId, action.ItemId)
			requiresTallyComputation = true

		case common.INCREMENT_TAKER_ID:
			s.logger.DebugLog(fmt.Sprintf("incrementing taker(takerID: %v)", action.TakerId))
			requiresTallyComputation = true

		case common.DECREMENT_TAKER_ID:
			s.logger.DebugLog(fmt.Sprintf("decrementing taker(takerID: %v)", action.TakerId))
			requiresTallyComputation = true

		// taker cases
		case common.ADD_NEW_TAKER:
			s.logger.DebugLog(fmt.Sprintf("adding taker(takerID: %v) with name : %v to split(splitID: %v)", action.TakerId, action.ItemName, splitID))
			err = s.ModelHelper.AddNewTakerToSplit(splitID, action.TakerId, action.ItemName)
			requiresTallyComputation = true

		case common.DELETE_TAKER:
			s.logger.DebugLog(fmt.Sprintf("deleting taker(takerID: %v) from split(splitID: %v)", action.TakerId, splitID))
			err = s.ModelHelper.DeleteTakerFromSplit(splitID, action.TakerId)
			requiresTallyComputation = true

		case common.EDIT_TAKER:
			s.logger.DebugLog(fmt.Sprintf("editing taker(takerID: %v)", action.TakerId))
			requiresTallyComputation = true

		case common.EDIT_SPLIT_TOTAL:
			s.logger.DebugLog(fmt.Sprintf("editing split total in split(splitID: %v) to %v", splitID, action.Total))
			err = s.ModelHelper.UpdateSplitTotal(splitID, action.Total)
			requiresTallyComputation = true
		}
	}

	if err == nil {
		if sess, ok := s.getSession(splitID); ok {
			sess.RequiresNewTally.Store(requiresTallyComputation)
			sess.BroadcastChannel <- common.ActionExecutionSuccessResponse(&action)
		}
	} else {
		s.logger.DebugLog(fmt.Sprintf("Failed to execute action: %v", err))
	}

	return err
}
