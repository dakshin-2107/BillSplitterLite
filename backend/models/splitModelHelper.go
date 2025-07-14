package models

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

type SplitModelHelper struct {
	db     IDatabase
	logger logger.ILogger
}

func (sp *SplitModelHelper) Init(db IDatabase, logger logger.ILogger) {
	sp.db = db
	sp.logger = logger
}

func (sp *SplitModelHelper) AddBillIfNotExists(split *Split) (string, error) {
	// valdiate split here if needed

	return sp.db.AddBillIfNotExists(split)
}

func (sp *SplitModelHelper) DeleteBillIfExists(splitId string) error {
	// delete split only once there is no activity from any participants

	return sp.db.DeleteBillIfExists(splitId)
}

func (sp *SplitModelHelper) UpdateBillMetaData(splitId string, newSplit Split) error {

	return sp.db.UpdateBillMetaData(splitId, newSplit)
}

func (sp *SplitModelHelper) AddBillItem(splitId string, item Item) error {

	return sp.db.AddBillItem(splitId, item)
}

func (sp *SplitModelHelper) DeleteBillItem(splitId string, itemId string) error {

	return sp.db.DeleteBillItem(splitId, itemId)
}

func (sp *SplitModelHelper) AddItemTaker(splitId string, itemId string, takerId string) error {

	return sp.db.AddItemTaker(splitId, itemId, takerId)
}

func (sp *SplitModelHelper) DeleteItemTaker(splitId string, itemId string, takerId string) error {

	return sp.db.DeleteItemTaker(splitId, itemId, takerId)
}

// func CreateNewSplit() Split {
// 	return Split{
// 		BillID:       "",
// 		Items:        []Item{},
// 		Participants: make(map[string]string),
// 		Status:       BILL_STATUS_ACTIVE,
// 		Location:     "",
// 		Date:         time.Now(),
// 		CreatedAt:    time.Now(),
// 		LastUpdated:  time.Now(),
// 	}
// }
