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
			if !sess.RequiresNewTally.Load() {
				continue
			}
			// Get the current split (containing all bills)
			split, err := sess.ModelHelper.GetSplit(sess.SplitID)
			if err != nil {
				continue
			}

			// Initialize Tally
			tally := common.Tally{
				UserShares:  make(map[string]common.UserShare),
				BillNameMap: make(map[int]string),
				ActualTotal: 0,
			}

			// Accumulate across all bills
			for _, bill := range split.Bills {
				tally.ActualTotal += bill.TotalAmount
				currBillId := bill.BillID
				tally.BillNameMap[currBillId] = bill.Location + " - " + bill.Date

				// For each item, divide price by number of takers
				for _, item := range bill.Items {
					totalSharesInItem := 0
					for _, sharesCount := range item.Takers {
						totalSharesInItem += sharesCount
					}

					if totalSharesInItem == 0 {
						continue
					}

					sharePrice := item.Price / float32(totalSharesInItem)

					for takerId, sharesCount := range item.Takers {
						shareForThisItem := sharePrice * float32(sharesCount)

						userShare, ok := tally.UserShares[takerId]
						if !ok {
							userShare = common.UserShare{
								BillShares:     make(map[int]common.BillShare),
								UserShareTotal: 0.0,
							}
						}

						billShare, ok := userShare.BillShares[currBillId]
						if !ok {
							billShare = common.BillShare{
								ItemShares:     make(map[string]float32),
								BillShareTotal: 0.0,
							}
						}

						billShare.ItemShares[item.Name] = shareForThisItem
						billShare.BillShareTotal += shareForThisItem
						userShare.BillShares[currBillId] = billShare

						userShare.UserShareTotal += shareForThisItem
						tally.UserShares[takerId] = userShare
					}
				}
			}

			// Calculate total and difference
			tally.CalculatedTotal = 0.0
			for _, us := range tally.UserShares {
				tally.CalculatedTotal += us.UserShareTotal
			}

			if split.TotalAmount > 0 {
				tally.ActualTotal = split.TotalAmount
			}

			tally.TotalDifference = tally.CalculatedTotal - tally.ActualTotal

			// Publish to everyone
			fmt.Println("Publishing tally: ", tally)
			sess.RequiresNewTally.Store(false)
			sess.BroadcastChannel <- common.TallyResponse(tally)
		}
	}
}
