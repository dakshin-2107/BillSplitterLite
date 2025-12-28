package handlers

import (
	"fmt"
	"net/http"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

type SplitActionWsHandler struct{ common.SessionRouteHandlerBase }

func (sp *SplitActionWsHandler) Init() {
	sp.Pattern = `/ws-actions`
	sp.Method = "GET"
	sp.Group = ""
	sp.Handler = sp.Handle
}

/*
- Since the initial request to this websocket endpoint is a GET request, extract session (and user) id from cookies
- Upgrade the connection to a websocket connection
- All user actions are sent as JSON messages through the websocket connection
*/
func (sp *SplitActionWsHandler) Handle(ctx *gin.Context) {

	// Extract bill id and user id from cookies
	billID, err := ctx.Cookie(common.BILL_SESSION_ID_STR)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.MissingCookieResponse())
		return
	}

	userID, err := ctx.Cookie(common.USER_SESSION_ID_STR)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.MissingCookieResponse())
		return
	}

	conn, err := sp.SessionManager.GetUserSessionConnection(billID, userID, ctx)
	if err != nil {
		sp.Logger.DebugLog(fmt.Sprintf("WebSocket upgrade failed: %v", err))
		return
	}

	defer conn.Close()
	sp.Logger.DebugLog("Established websocket connection")
	for {

		_, actionMsg, err := conn.ReadMessage()
		if err != nil {
			conn, err = sp.SessionManager.GetUserSessionConnection(billID, userID, ctx)
			if err != nil {
				sp.Logger.DebugLog(fmt.Sprintf("WebSocket upgrade failed: %v", err))
				return
			}
		}

		sp.SessionManager.ExecuteAction(billID, userID, actionMsg)
	}
}

func testSplitHandlerWs() common.ISessionRouteHandlerBase {
	return &SplitActionWsHandler{}
}

func testSplitHandler2Ws() common.IRouteHandlerBase {
	return &SplitActionWsHandler{}
}
