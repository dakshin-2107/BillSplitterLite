package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

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
	ctx.SetCookie(common.BILL_SESSION_ID_STR, "", -1, "/", "localhost", true, true)
	ctx.SetCookie(common.USER_SESSION_ID_STR, "", -1, "/", "localhost", true, true)
	ctx.SetCookie(common.HELLO_THERE_STR, "", -1, "/", "localhost", false, false)
}
