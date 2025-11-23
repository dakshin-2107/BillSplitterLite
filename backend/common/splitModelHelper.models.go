package common

import (
	"strings"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

const BILL_STATUS_COMPLETED = 1
const BILL_STATUS_ACTIVE = 0

type Split struct {
	BillID       string            `json:"billId"`
	Items        map[string]Item   `json:"items"`
	TotalAmount  float64           `json:"total"`
	Participants map[string]string `json:"participants"`
	Status       int               `json:"status"`
	Location     string            `json:"location"`
	Date         SplitDate         `json:"date"`
	CreatedAt    SplitDate         `json:"createdAt"`
	LastUpdated  SplitDate         `json:"lastUpdated"`
}

type Item struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Price  float32        `json:"price"`
	Takers map[string]int `json:"takers"`
}

type ISplitModifier interface {
	GetBill(string) (*Split, error)
	AddBillIfNotExists(split *Split) error
	DeleteBillIfExists(splitId string) error

	UpdateBillMetaData(splitId string, newSplit Split) error
	AddBillItem(splitId string, item Item) error
	DeleteBillItem(splitId string, itemId string) error

	AddItemTaker(splitId string, itemId string, takerId string) error
	DeleteItemTaker(splitId string, itemId string, takerId string) error

	AddTakerID(splitId string, takerId string, takerName string) error
	DeleteTakerID(splitId string, takerId string) error
}

type ISplitModelHelper interface {
	Init(db IDatabase, logger logger.ILogger)
	CreateNewSplit() *Split
	ISplitModifier
}

type SplitDate struct{ string }

// object -> JSON type
func (d SplitDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.string + `"`), nil
}

// JSON string -> object
func (d *SplitDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("02/01/2006", s)
	if err == nil {
		d.string = t.Format("02/01/2006")
		return nil
	}

	t, err = time.Parse("02/01/06", s)
	if err == nil {
		d.string = t.Format("02/01/2006")
		return nil
	}

	t, err = time.Parse("02-Jan-2006", s)
	if err == nil {
		d.string = t.Format("02/01/2006")
		return nil
	}

	t, err = time.Parse("02-Jan-06", s)
	if err == nil {
		d.string = t.Format("02/01/2006")
		return nil
	}

	return err
}

func CreateSplitDate(t time.Time) SplitDate {
	// dateTimeStr := t.Format("02/01/2006 15:04:05")
	dateStr := t.Format("02/01/2006")
	return SplitDate{string: dateStr}
}
