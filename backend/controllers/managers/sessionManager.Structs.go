package managers

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Action struct {
	ActionId   int     `json:"actionId"`
	ActionType int     `json:"actionType"`
	SplitId    string  `json:"splitId"`
	BillId     int     `json:"billId"`
	ItemId     int     `json:"itemId"`
	ItemName   string  `json:"itemName"` // can be used for taker name
	TakerId    string  `json:"takerId"`
	Price      float32 `json:"price"`
	Total      float32 `json:"total"`
	Location   string  `json:"location"`
	Date       string  `json:"date"`
}

type SessionManager struct {
	mu          sync.RWMutex // protects SessionDict
	SessionDict map[string]*SplitSession
	Upgrader    *websocket.Upgrader
	logger      logger.ILogger
	ModelHelper common.ISplitModelHelper
}

type SplitSession struct {
	SplitID           string
	ClientConnections map[string]*websocket.Conn
	AdminConnection   *websocket.Conn
	UserIdCounter     int
	RequiresNewTally  atomic.Bool
	LastUsed          time.Time
	BroadcastChannel  chan gin.H
	SignalChannel     chan Action
	ModelHelper       common.ISplitModelHelper
}

// dummy method for now, will remove it if not needed.
func (a *Action) ExecuteAction() error {
	return nil
}
