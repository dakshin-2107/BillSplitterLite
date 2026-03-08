package handlers

import (
	"net/http"
	"os"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

/*
Purpose:
- deletes the cookies that represent the users session
*/
type CloseHandler struct {
	common.SessionRouteHandlerBase
}

func (ch *CloseHandler) Init() {
	ch.Pattern = "/close"
	ch.Method = "GET"
	ch.Group = ""
	ch.Handler = ch.Handle
}

func (ch *CloseHandler) Handle(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	domain := os.Getenv("BASE_URL")
	ctx.SetCookie(common.BILL_SESSION_ID_STR, "", -1, "/", domain, false, true)
	ctx.SetCookie(common.USER_SESSION_ID_STR, "", -1, "/", domain, false, true)
}
