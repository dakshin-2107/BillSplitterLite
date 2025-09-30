package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/gin-gonic/gin"
)

// Endpoint to just check whether the server is up and running, does not need anything
type PingHandler struct{ RouteHandlerBase }

func (p *PingHandler) Init(logger logger.ILogger, modelHelper models.ISplitModelHelper, sessionManger managers.ISessionManager) {
	p.RouteHandlerBase.BaseInit(logger, modelHelper, sessionManger)

	p.pattern = `/ping`
	p.method = "GET"
	p.group = ""
	p.handler = p.Handle
}

func (p *PingHandler) Handle(ctx *gin.Context) {
	p.logger.DebugLog("Server was just pinged")
	ctx.JSON(200, "pong")
}

func testPing() IRouteHandlerBase {
	return &PingHandler{}
}
