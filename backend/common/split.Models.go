package common

import (
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

type Bill struct {
	BillID        int             `json:"billId"`
	Items         map[string]Item `json:"items"`
	TotalAmount   float32         `json:"total"`
	ItemIdCounter int             `json:"itemIdCounter"`
	Location      string          `json:"location"`
	Date          string          `json:"date"`
}

type Split struct {
	SplitID       string            `json:"splitId"`
	Bills         map[int]Bill      `json:"bills"`
	TotalAmount   float32           `json:"totalAmount"`
	Participants  map[string]string `json:"participants"`
	BillIdCounter int               `json:"billIdCounter"`
	CreatedAt     time.Time         `json:"createdAt"`
	LastUpdated   time.Time         `json:"lastUpdated"`
}

type Item struct {
	Id     int            `json:"id"`
	Name   string         `json:"name"`
	Price  float32        `json:"price"`
	Takers map[string]int `json:"takers"`
}

type ISplitModifier interface {
	GetSplit(splitId string) (*Split, error)
	GetNewBillId(splitId string) int

	AddSplitIfNotExists(split *Split) error
	DeleteSplitIfExists(splitId string) error

	AddNewTakerToSplit(splitId string, takerId string, takerName string) error
	DeleteTakerFromSplit(splitId string, takerId string) error

	AddBillToSplit(splitId string, bill Bill) error
	DeleteBillFromSplit(splitId string, billId int) error

	GetNewItemId(splitId string, billId int) int
	AddItemToBill(splitId string, billId int, item Item) error
	DeleteItemFromBill(splitId string, billId int, itemId int) error

	AddNewTakerForItem(splitId string, billId int, itemId int, takerId string) error
	DeleteTakerForItem(splitId string, billId int, itemId int, takerId string) error

	UpdateBillInformation(splitId string, billId int, newTotal float32, newLocation string, newDate string) error
}

type ISplitModelHelper interface {
	Init(db IDatabase, logger logger.ILogger)
	CreateNewSplit() *Split
	ISplitModifier
}
