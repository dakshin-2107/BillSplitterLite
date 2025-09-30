package managers

import (
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gorilla/websocket"
)

// Action type
const (
	HELLO_THERE = iota
	ADD_TAKER
	REMOVE_TAKER
	ADD_ITEM
	REMOVE_ITEM
	EDIT_ITEM
)

type ISessionManager interface {
	Init(logger logger.ILogger, connUpgrader *ConnectionUpgrader)
	GetSessionConnnections(sessionId string) ([]*MyWebSocketConnection, error)
	CreateNewSession(sessionId string) error
	AddConnection(sessionId string, newConnection *MyWebSocketConnection) error
	CleanSessions() error
	DeleteSession(sessionId string) error
	IActionService
}

type IActionManager interface {
	QueueAction(action IAction) error
	DequeueAction() (IAction, error)
}

type IActionService interface {
	RunActionService()
	PublishAction(action IAction) error
	IActionManager
}

type IAction interface {
	ExecuteAction() error
}

// Example JSON:
//
//	{
//	  "ActionType": 1,
//	  "SplitId": "split123",
//	  "ItemId": "item456",
//	  "ItemName": "Pizza",
//	  "TakerId": "user789",
//	  "Price": 299.99
//	}
type Action struct {
	ActionType int     `redis:"actionType"`
	SplitId    string  `redis:"splitId"`
	ItemId     string  `redis:"itemId"`
	ItemName   string  `redis:"itemName"`
	TakerId    string  `redis:"takerId"`
	Price      float32 `redis:"price"`
}

type SessionManager struct {
	ConnectionDict map[string]*SplitSession
	upgrader       *ConnectionUpgrader
	logger         logger.ILogger
}

type SplitSession struct {
	ClientConnections []*MyWebSocketConnection
	AdminConnection   *MyWebSocketConnection
	IsSessionActive   bool
	LastUsed          time.Time
}

type MyWebSocketConnection struct{ websocket.Conn }

type ConnectionUpgrader struct{ websocket.Upgrader }

func testManager() ISessionManager {
	return &SessionManager{}
}
