package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

/*
Purpose :
- Used only for pinging
*/

type PingHandler struct{ common.RouteHandlerBase }

func (p *PingHandler) Init() {
	p.Pattern = `/ping`
	p.Method = "GET"
	p.Group = ""
	p.Handler = p.Handle
}

func (p *PingHandler) Handle(ctx *gin.Context) {
	p.Logger.DebugLog("Server was just pinged")
	ctx.JSON(200, "pong")
}

func testPing() common.IRouteHandlerBase {
	return &PingHandler{}
}
