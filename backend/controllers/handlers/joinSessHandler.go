package handlers

import (
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

type SharedSessionHandler struct {
	common.SessionRouteHandlerBase
}

func (sh *SharedSessionHandler) Init() {
	sh.Pattern = `/join/:billId`
	sh.Method = "GET"
	sh.Group = ""
	sh.Handler = sh.Handle
}

func (sh *SharedSessionHandler) Handle(ctx *gin.Context) {
	billID := ctx.Param("billId")
	if billID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bill ID is missing",
		})
		return
	}

	userID, err := sh.SessionManager.GetNewUserSessionId(billID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get new user session ID",
		})
		return
	}

	domain := os.Getenv("BASE_URL")
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(common.BILL_SESSION_ID_STR, billID, math.MaxInt64, "/", domain, false, true)
	ctx.SetCookie(common.USER_SESSION_ID_STR, fmt.Sprintf("user_%s", userID), math.MaxInt64, "/", domain, false, true)
	ctx.JSON(http.StatusOK, common.SessionUriResponse(billID))
}
