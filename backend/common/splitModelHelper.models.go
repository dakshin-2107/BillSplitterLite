package common

import (
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

type Split struct {
	BillID        string            `json:"billId"`
	Items         map[string]Item   `json:"items"`
	TotalAmount   float32           `json:"total"`
	Participants  map[string]string `json:"participants"` // map of takerName -> takerId
	ItemIdCounter int               `json:"itemIdCounter"`
	Location      string            `json:"location"`
	Date          string            `json:"date"`
	CreatedAt     time.Time         `json:"createdAt"`
	LastUpdated   time.Time         `json:"lastUpdated"`
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
	GetNewItemId(splitID string) string
	AddBillItem(splitId string, item Item) error
	DeleteBillItem(splitId string, itemId string) error

	AddItemTaker(splitId string, itemId string, takerId string) error
	DeleteItemTaker(splitId string, itemId string, takerId string) error

	AddTakerID(splitId string, takerId string, takerName string) error
	DeleteTakerID(splitId string, takerId string) error

	UpdateBillInformation(splitId string, newTotal float32) error
}

type ISplitModelHelper interface {
	Init(db IDatabase, logger logger.ILogger)
	CreateNewSplit() *Split
	ISplitModifier
}
