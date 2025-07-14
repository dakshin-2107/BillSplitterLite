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
	Participants map[string]string `json:"participants"`
	Status       int               `json:"status"`
	Location     string            `json:"location"`
	Date         time.Time         `json:"date"`
	CreatedAt    time.Time         `json:"createdAt"`
	LastUpdated  time.Time         `json:"lastUpdated"`
}

type Item struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Price  float64  `json:"price"`
	Takers []string `json:"takers"`
}

type ISplitModifier interface {
	AddBillIfNotExists(split *Split) (string, error)
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
	ISplitModifier
}

/*

JSON schema to be followed

{
  "billId": "unique-bill-link-id",
  "items": [
    {
      "id": "item1",
      "name": "Pizza",
      "price": 25.00,
      "takers": ["userA", "userB"]
    },
    {
      "id": "item2",
      "name": "Coke",
      "price": 3.00,
      "takers": ["userC"]
    }
  ],
  "participants": {
    "userA": "Alice",
    "userB": "Bob",
    "userC": "Charlie",
  },
  "status": "active",
  "location": "place name",
  "date": "date",
  "createdAt": "...",
  "lastUpdated": "..."
}

*/
