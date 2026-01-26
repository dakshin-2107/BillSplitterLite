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

	// edit bill information - actual bill total for now
	EDIT_SPLIT_TOTAL

	// add all participants to an item
	ADD_ALL_TAKERS_FOR_ITEM

	// close the session and delete all the data associated with it
	BYE_BYE
)

type ISessionManager interface {
	Init(logger logger.ILogger, connUpgrader *websocket.Upgrader, modelHelper ISplitModelHelper)
	GetNewUserSessionId(sessionId string) (string, error)
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

/*
Example tally JSON
{
	userShares: {
		Kevin : {
			billShares : {
				"1" : {
					itemShares : {
						"1" : 100.00,
						"2" : 200.00
					},
					billShareTotal : 300.00
				},
				"2" : {
					itemShares : {
						"1" : 150.00,
						"2" : 200.00
					},
					billShareTotal : 350.00
				}
			},
			userShareTotal : 650.00
		}
	},
	billNameMap: {
		"1" : "Poor richards - 1st Jan, 2025",
		"2" : "Pizza by alfredo's - 25th Dec, 2025"
	},
	actualTotal: 1000.00,
	calculatedTotal: 1000.00,
	totalDifference: 0.00
}

*/

type Tally struct {
	UserShares      map[string]UserShare `json:"userShares"`
	BillNameMap     map[int]string       `json:"billNameMap"`
	ActualTotal     float32              `json:"actualTotal"`
	CalculatedTotal float32              `json:"calculatedTotal"`
	TotalDifference float32              `json:"totalDifference"`
}

type UserShare struct {
	BillShares     map[int]BillShare `json:"billShares"`
	UserShareTotal float32           `json:"userShareTotal"`
}

type BillShare struct {
	ItemShares     map[string]float32 `json:"itemShares"`
	BillShareTotal float32            `json:"billShareTotal"`
}
