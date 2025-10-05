package modelHelper

import (
	"fmt"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

// TODO: Add all domain specific validations before the DB is called
type SplitModelHelper struct {
	db     common.IDatabase
	logger logger.ILogger
}

func testModel() common.ISplitModelHelper {
	return &SplitModelHelper{}
}

func (sp *SplitModelHelper) Init(db common.IDatabase, logger logger.ILogger) {
	sp.db = db
	sp.logger = logger
}

// Bill methods
func (sp *SplitModelHelper) AddBillIfNotExists(split *common.Split) error {
	return sp.db.AddBillIfNotExists(split)
}

func (sp *SplitModelHelper) DeleteBillIfExists(splitId string) error {
	return sp.db.DeleteBillIfExists(splitId)
}

func (sp *SplitModelHelper) UpdateBillMetaData(splitId string, newSplit common.Split) error {
	return sp.db.UpdateBillMetaData(splitId, newSplit)
}

// Bill item methods
func (sp *SplitModelHelper) AddBillItem(splitId string, item common.Item) error {
	return sp.db.AddBillItem(splitId, item)
}

func (sp *SplitModelHelper) DeleteBillItem(splitId string, itemId string) error {
	return sp.db.DeleteBillItem(splitId, itemId)
}

// Taker methods
func (sp *SplitModelHelper) AddItemTaker(splitId string, itemId string, takerId string) error {
	sp.logger.DebugLog(fmt.Sprintf("Adding item taker %v for item %v for split ID %v", takerId, itemId, splitId))
	return sp.db.AddItemTaker(splitId, itemId, takerId)
}

func (sp *SplitModelHelper) DeleteItemTaker(splitId string, itemId string, takerId string) error {
	return sp.db.DeleteItemTaker(splitId, itemId, takerId)
}

func (sp *SplitModelHelper) GetBill(billId string) (*common.Split, error) {
	return sp.db.GetBill(billId)
}

func (sp *SplitModelHelper) CreateNewSplit() *common.Split {
	return &common.Split{
		BillID:       "",
		Items:        make(map[string]common.Item),
		TotalAmount:  0.0,
		Participants: make(map[string]string),
		Status:       common.BILL_STATUS_ACTIVE,
		Location:     "",
		Date:         common.CreateSplitDate(time.Now()),
		CreatedAt:    common.CreateSplitDate(time.Now()),
		LastUpdated:  common.CreateSplitDate(time.Now()),
	}
}
