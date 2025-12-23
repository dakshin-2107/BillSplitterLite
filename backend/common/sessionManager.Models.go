package common

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Action types
const (
	// initiate the session
	HELLO_THERE = iota

	// just to the ping the backend
	PING

	// sync the current bill state with the client
	SYNC_BILL_STATE

	// basic taker CRUD operations
	ADD_TAKER_FOR_ITEM // operation must be idempotent
	DELETE_TAKER_FOR_ITEM

	// increment/decrement the number of shares for a taker in an item
	INCREMENT_TAKER_ID
	DECREMENT_TAKER_ID

	// basic item CRUD operations
	ADD_NEW_ITEM
	DELETE_ITEM
	EDIT_ITEM

	// create/remove the taker only from the item
	ADD_NEW_TAKER
	DELETE_TAKER
	EDIT_TAKER

	// close the session and delete all the data associated with it
	BYE_BYE
)

type ISessionManager interface {
	Init(logger logger.ILogger, connUpgrader *websocket.Upgrader, modelHelper ISplitModelHelper)
	GetSessionConnnections(sessionId string) (map[string]*websocket.Conn, error)
	CreateNewSession(sessionId string, userid string, ctx *gin.Context) (*websocket.Conn, error)
	AddConnection(sessionId string, newConnection *websocket.Conn) error
	CleanSessions() error
	DeleteSession(sessionId string) error
	IActionService
}

type IActionManager interface {
	QueueAction(action IAction) error
	DequeueAction() (IAction, error)
}

type IActionService interface {
	ExecuteAction(string, string, []byte) (IAction, error)
	PublishAction(action IAction) error
	IActionManager
}

type IAction interface {
	ExecuteAction() error
}

// type MyWebSocketConnection struct{ websocket.Conn }

// type ConnectionUpgrader struct{ websocket.Upgrader }
