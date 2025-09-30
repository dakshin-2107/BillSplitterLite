package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

type SplitModelHelper struct {
	db     IDatabase
	logger logger.ILogger
}

func testModel() ISplitModelHelper {
	return &SplitModelHelper{}
}

func (sp *SplitModelHelper) Init(db IDatabase, logger logger.ILogger) {
	sp.db = db
	sp.logger = logger
}

// Bill methods
func (sp *SplitModelHelper) AddBillIfNotExists(split *Split) error {

	return sp.db.AddBillIfNotExists(split)
}

func (sp *SplitModelHelper) DeleteBillIfExists(splitId string) error {

	return sp.db.DeleteBillIfExists(splitId)
}

func (sp *SplitModelHelper) UpdateBillMetaData(splitId string, newSplit Split) error {

	return sp.db.UpdateBillMetaData(splitId, newSplit)
}

// Bill item methods
func (sp *SplitModelHelper) AddBillItem(splitId string, item Item) error {

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

func (sp *SplitModelHelper) CreateNewSplit() *Split {
	return &Split{
		BillID:       "",
		Items:        []Item{},
		TotalAmount:  0.0,
		Participants: make(map[string]string),
		Status:       BILL_STATUS_ACTIVE,
		Location:     "",
		Date:         sp.CreateSplitDate(time.Now()),
		CreatedAt:    sp.CreateSplitDate(time.Now()),
		LastUpdated:  sp.CreateSplitDate(time.Now()),
	}
}

// string -> JSON
func (d SplitDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.string + `"`), nil
}

// JSON -> string
func (d *SplitDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Try to parse the date in DD/MM/YYYY format
	t, err := time.Parse("02/01/2006", s)
	if err == nil {
		d.string = t.Format("02/01/2006")
		return nil
	}

	t, err = time.Parse("02/01/06", s)
	if err == nil {
		d.string = t.Format("02/01/2006")
	}

	return err
}

func (sp *SplitModelHelper) CreateSplitDate(t time.Time) SplitDate {
	// dateTimeStr := t.Format("02/01/2006 15:04:05")
	dateStr := t.Format("02/01/2006")
	return SplitDate{string: dateStr}
}
