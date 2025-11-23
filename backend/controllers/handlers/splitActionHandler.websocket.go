package handlers

import (
	"fmt"
	"net/http"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

/*
Purpose :
  - use the bill id and user id from a cookie
  - use long polling to get updates from the user on need by basis
  - use the session manager to execute actions and then publish the confirmed updates to the user
*/

type SplitActionWsHandler struct{ common.SessionRouteHandlerBase }

func (sp *SplitActionWsHandler) Init() {
	sp.Pattern = `/split`
	sp.Method = "GET"
	sp.Group = ""
	sp.Handler = sp.Handle
}

func (sp *SplitActionWsHandler) Handle(ctx *gin.Context) {

	// Extract bill id and user id from cookies
	billID, err := ctx.Cookie("BillsessionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing bill_id cookie"})
		return
	}

	userID, err := ctx.Cookie("UserId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id cookie"})
		return
	}

	conn, err := sp.SessionManager.CreateNewSession(billID, userID, ctx)
	if err != nil {
		sp.Logger.DebugLog(fmt.Sprintf("WebSocket upgrade failed: %v", err))
		return
	}

	defer conn.Close()

	for {
		// Read action from frontend
		_, msg, err := conn.ReadMessage()
		if err != nil {
			sp.Logger.DebugLog(fmt.Sprintf("Read error: %v", err))
			break
		}

		sp.Logger.DebugLog(string(msg))

		if action, err := sp.SessionManager.ExecuteAction(billID, userID, msg); err != nil {
			conn.WriteJSON(gin.H{
				"success": true,
				"action":  action,
			})
		}

		// do something  here

		conn.WriteJSON(gin.H{
			"success": false,
			"message": "Could not perform action",
		})

		// if err := conn.WriteJSON(ack); err != nil {
		// 	log.Printf("Write error: %v", err)
		// 	break
		// }
	}
}

func (sp *SplitActionWsHandler) ParseAction(message []byte) interface{} {

	return nil
}

func testSplitHandlerWs() common.ISessionRouteHandlerBase {
	return &SplitActionWsHandler{}
}

func testSplitHandler2Ws() common.IRouteHandlerBase {
	return &SplitActionWsHandler{}
}
