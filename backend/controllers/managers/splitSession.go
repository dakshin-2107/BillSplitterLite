package managers

import (
	"fmt"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
)

// TODO : if there is no new messages from Broadcast channel for more than 15mins, then we delete the session
func (s *SplitSession) RunPublisherService() {
	for {
		select {
		case action := <-s.SignalChannel:
			if action.ActionType == common.BYE_BYE {
				return
			}

		case chMessage := <-s.BroadcastChannel:
			for _, conn := range s.ClientConnections {
				conn.WriteJSON(chMessage)
			}
		}
	}
}

func (sess *SplitSession) RunTallyService() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case action := <-sess.SignalChannel:
			if action.ActionType == common.BYE_BYE {
				return
			}
		case <-ticker.C:
			if !sess.RequiresNewTally {
				continue
			}
			// Get the current bill
			bill, err := sess.ModelHelper.GetBill(sess.SplitID)
			if err != nil {
				continue
			}

			// Initialize Tally
			tally := common.Tally{
				UserShares:  make(map[string]common.UserShare),
				ActualTotal: float32(bill.TotalAmount),
			}

			// For each item, divide price by number of takers
			for _, item := range bill.Items {
				totalSharesInItem := 0
				for _, shares := range item.Takers {
					totalSharesInItem += shares
				}

				if totalSharesInItem == 0 {
					continue
				}

				sharePrice := item.Price / float32(totalSharesInItem)

				for takerId, sharesCount := range item.Takers {
					userShare, ok := tally.UserShares[takerId]
					if !ok {
						userShare = common.UserShare{
							Shares: make(map[string]float32),
						}
					}

					shareForThisItem := sharePrice * float32(sharesCount)
					userShare.Shares[item.Name] = shareForThisItem
					userShare.UserShareTotal += shareForThisItem
					tally.UserShares[takerId] = userShare
				}
			}

			// Calculate total and difference
			tally.CalculatedTotal = 0.0
			for _, us := range tally.UserShares {
				tally.CalculatedTotal += us.UserShareTotal
			}

			tally.TotalDifference = tally.CalculatedTotal - tally.ActualTotal

			// 5. Publish to everyone
			fmt.Println("Publishing tally: ", tally)
			sess.RequiresNewTally = false
			sess.BroadcastChannel <- common.TallyResponse(tally)
		}
	}
}
