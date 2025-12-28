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

	// edit bill information
	EDIT_BILL_INFORMATION

	// close the session and delete all the data associated with it
	BYE_BYE
)

type ISessionManager interface {
	Init(logger logger.ILogger, connUpgrader *websocket.Upgrader, modelHelper ISplitModelHelper)
	GetAllSessionConnections(sessionId string) (map[string]*websocket.Conn, error)
	GetUserSessionConnection(sessionId string, userid string, ctx *gin.Context) (*websocket.Conn, error)
	AddConnection(sessionId string, newConnection *websocket.Conn) error
	CleanSessions() error
	DeleteSession(sessionId string) error
	IActionService
}

type IActionService interface {
	ExecuteAction(string, string, []byte) error
}

type IAction interface {
	ExecuteAction() error
}

type IPublisher interface {
	Run()
	PublishAction(action IAction) error
}

type Tally struct {
	UserShares      map[string]UserShare `json:"userShares"`
	ActualTotal     float32              `json:"actualTotal"`
	CalculatedTotal float32              `json:"calculatedTotal"`
	TotalDifference float32              `json:"totalDifference"`
}

type UserShare struct {
	Shares         map[string]float32 `json:"shares"`
	UserShareTotal float32            `json:"userShareTotal"`
}
