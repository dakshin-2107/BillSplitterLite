package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

/*
Purpose :
  - use the bill id and user id from a cookie
  - upgrade the connection after getting those details
  - every action received from the front end must be validated, applied and then acknowledged back so that the frontend can be udpated
*/

type SplitActionHandler struct{ common.SessionRouteHandlerBase }

func (sp *SplitActionHandler) Init() {
	sp.Pattern = `/split`
	sp.Method = "POST"
	sp.Group = ""
	sp.Handler = sp.Handle
}

func (sp *SplitActionHandler) Handle(ctx *gin.Context) {

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

	var body struct {
		Action json.RawMessage `json:"action"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	msg := body.Action

	sp.Logger.DebugLog(string(msg))

	if action, err := sp.SessionManager.ExecuteAction(billID, userID, msg); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "action was executed",
			"action":  action,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "action could not be executed",
			"action":  action,
			"err":     err,
		})
	}
}

func (sp *SplitActionHandler) ParseAction(message []byte) interface{} {

	return nil
}

func testSplitHandler() common.ISessionRouteHandlerBase {
	return &SplitActionHandler{}
}

func testSplitHandler2() common.IRouteHandlerBase {
	return &SplitActionHandler{}
}
