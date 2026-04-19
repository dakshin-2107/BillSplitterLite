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

func computeBillShare(bill common.Bill) common.BillShare {
	billShare := common.BillShare{
		BillName:   bill.Location + " - " + bill.Date,
		UserShares: make(map[string]common.UserShare),
		BillTotal:  0,
	}

	for _, item := range bill.Items {
		totalShares := 0
		for _, count := range item.Takers {
			totalShares += count
		}
		if totalShares == 0 {
			continue
		}

		sharePrice := item.Price / float32(totalShares)
		for takerId, count := range item.Takers {
			share := sharePrice * float32(count)

			us := billShare.UserShares[takerId]
			if us.ItemShares == nil {
				us.ItemShares = make(map[string]float32)
			}
			us.ItemShares[item.Name] = share
			us.UserShareTotal += share
			billShare.UserShares[takerId] = us

			billShare.BillTotal += share
		}
	}

	return billShare
}

func recalcTotal(shares map[int]common.BillShare) float32 {
	var total float32
	for _, bs := range shares {
		total += bs.BillTotal
	}
	return total
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
			requiresFull := sess.RequiresNewTally.Load()
			dirtyBills := sess.GetAndClearDirtyBills()

			if !requiresFull && len(dirtyBills) == 0 {
				continue
			}

			if requiresFull {
				split, err := sess.ModelHelper.GetSplit(sess.SplitID)
				if err != nil {
					continue
				}

				tally := common.Tally{BillShares: make(map[int]common.BillShare)}
				for _, bill := range split.Bills {
					tally.BillShares[bill.BillID] = computeBillShare(bill)
				}
				tally.CalculatedTotal = recalcTotal(tally.BillShares)

				sess.RequiresNewTally.Store(false)
				fmt.Println("Publishing full tally")
				sess.BroadcastChannel <- common.TallyResponse(tally)

			} else {
				// Partial recompute — fetch only dirty bills, let frontend merge into its tally cache
				dirtyBillShares := make(map[int]common.BillShare)
				for _, billId := range dirtyBills {
					bill, err := sess.ModelHelper.GetBill(sess.SplitID, billId)
					if err != nil {
						continue
					}
					dirtyBillShares[billId] = computeBillShare(*bill)
				}
				if len(dirtyBillShares) > 0 {
					fmt.Printf("Publishing partial tally (dirty bills: %v)\n", dirtyBills)
					sess.BroadcastChannel <- common.BillTallyResponse(dirtyBillShares)
				}
			}
		}
	}
}
