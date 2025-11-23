package common

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Action type
const (
	HELLO_THERE       = iota
	ADD_ITEM_TAKER    // 1
	REMOVE_ITEM_TAKER // 2
	ADD_ITEM          // 3
	REMOVE_ITEM       // 4
	EDIT_ITEM         // 5
	ADD_TAKER_ID      // 6
	REMOVE_TAKER_ID   // 7
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
