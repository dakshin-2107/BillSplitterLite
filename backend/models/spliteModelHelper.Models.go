package models

import (
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

const BILL_STATUS_COMPLETED = 1
const BILL_STATUS_ACTIVE = 0

type Split struct {
	BillID       string            `json:"billId"`
	Items        []Item            `json:"items"`
	TotalAmount  float64           `json:"total"`
	Participants map[string]string `json:"participants"`
	Status       int               `json:"status"`
	Location     string            `json:"location"`
	Date         SplitDate         `json:"date"`
	CreatedAt    SplitDate         `json:"createdAt"`
	LastUpdated  SplitDate         `json:"lastUpdated"`
}

type Item struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Price  float64  `json:"price"`
	Takers []string `json:"takers"`
}

type ISplitModifier interface {
	AddBillIfNotExists(split *Split) error
	DeleteBillIfExists(splitId string) error

	// admin actions
	UpdateBillMetaData(splitId string, newSplit Split) error
	AddBillItem(splitId string, item Item) error
	DeleteBillItem(splitId string, itemId string) error

	// admin + other user actions
	AddItemTaker(splitId string, itemId string, takerId string) error
	DeleteItemTaker(splitId string, itemId string, takerId string) error
}

type ISplitModelHelper interface {
	Init(db IDatabase, logger logger.ILogger)
	CreateNewSplit() *Split
	CreateSplitDate(time.Time) SplitDate
	ISplitModifier
}

type SplitDate struct {
	string string
}
