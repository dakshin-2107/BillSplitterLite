package modelHelper

import (
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

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

// Split management
func (sp *SplitModelHelper) GetSplit(splitId string) (*common.Split, error) {
	return sp.db.GetSplit(splitId)
}

func (sp *SplitModelHelper) AddSplitIfNotExists(split *common.Split) error {
	return sp.db.AddSplitIfNotExists(split)
}

func (sp *SplitModelHelper) DeleteSplitIfExists(splitId string) error {
	return sp.db.DeleteSplitIfExists(splitId)
}

func (sp *SplitModelHelper) GetNewBillId(splitId string) int {
	return sp.db.GetNewBillId(splitId)
}

// Taker management
func (sp *SplitModelHelper) AddNewTakerToSplit(splitId string, takerId string, takerName string) error {
	return sp.db.AddNewTakerToSplit(splitId, takerId, takerName)
}

func (sp *SplitModelHelper) DeleteTakerFromSplit(splitId string, takerId string) error {
	return sp.db.DeleteTakerFromSplit(splitId, takerId)
}

// Bill management
func (sp *SplitModelHelper) AddBillToSplit(splitId string, bill common.Bill) error {
	return sp.db.AddBillToSplit(splitId, bill)
}

func (sp *SplitModelHelper) DeleteBillFromSplit(splitId string, billId int) error {
	return sp.db.DeleteBillFromSplit(splitId, billId)
}

func (sp *SplitModelHelper) UpdateBillInformation(splitId string, billId int, newTotal float32, newLocation string, newDate string) error {
	return sp.db.UpdateBillInformation(splitId, billId, newTotal, newLocation, newDate)
}

// Item management
func (sp *SplitModelHelper) GetNewItemId(splitId string, billId int) int {
	return sp.db.GetNewItemId(splitId, billId)
}

func (sp *SplitModelHelper) AddItemToBill(splitId string, billId int, item common.Item) error {
	return sp.db.AddItemToBill(splitId, billId, item)
}

func (sp *SplitModelHelper) DeleteItemFromBill(splitId string, billId int, itemId int) error {
	return sp.db.DeleteItemFromBill(splitId, billId, itemId)
}

func (sp *SplitModelHelper) AddNewTakerForItem(splitId string, billId int, itemId int, takerId string) error {
	return sp.db.AddNewTakerForItem(splitId, billId, itemId, takerId)
}

func (sp *SplitModelHelper) DeleteTakerForItem(splitId string, billId int, itemId int, takerId string) error {
	return sp.db.DeleteTakerForItem(splitId, billId, itemId, takerId)
}

func (sp *SplitModelHelper) AddAllTakersToItem(splitId string, billId int, itemId int) error {
	return sp.db.AddAllTakersToItem(splitId, billId, itemId)
}

// Creation helpers
func (sp *SplitModelHelper) CreateNewSplit() *common.Split {
	return &common.Split{
		SplitID:       "",
		Bills:         make(map[int]common.Bill),
		TotalAmount:   0.0,
		Participants:  make(map[string]string),
		BillIdCounter: 0,
		CreatedAt:     time.Now(),
		LastUpdated:   time.Now(),
	}
}
