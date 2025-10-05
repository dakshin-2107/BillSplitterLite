package managers

import (
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gorilla/websocket"
)

// // Action type
// const (
// 	HELLO_THERE = iota
// 	ADD_TAKER = 1
// 	REMOVE_TAKER = 2
// 	ADD_ITEM = 3
// 	REMOVE_ITEM
// 	EDIT_ITEM
// )

/*
{
	"actionId" : 2,
	"actionType" : ADD_ITEM,
	"spliteId" : "asdasdasf34f4tr",
	"itemId" : "asdfkl3121",
	"itemName" : "roll",
	"takerId" : "1",
	"price" : 90.00
}
*/

type Action struct {
	ActionId   int     `json:"actionId"`
	ActionType int     `json:"actionType"`
	SplitId    string  `json:"splitId"`
	ItemId     string  `json:"itemId"`
	ItemName   string  `json:"itemName"`
	TakerId    string  `json:"takerId"`
	Price      float32 `json:"price"`
}

type SessionManager struct {
	ConnectionDict map[string]*SplitSession
	Upgrader       *websocket.Upgrader
	logger         logger.ILogger
	ModelHelper    common.ISplitModelHelper
}

type SplitSession struct {
	ClientConnections map[string]*websocket.Conn
	AdminConnection   *websocket.Conn
	ActionCount       int
	IsSessionActive   bool
	LastUsed          time.Time
}
