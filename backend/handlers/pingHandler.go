package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/gin-gonic/gin"
)

func testPing() IRouteHandlerBase {
	return &PingHandler{}
}

type PingHandler struct{ RouteHandlerBase }

func (p *PingHandler) Init(logger logger.ILogger, modelHelper models.ISplitModelHelper) {
	p.pattern = `/ping`
	p.method = "GET"
	p.group = ""
	p.logger = logger
	p.modelHelper = modelHelper
	p.handler = p.Handle
}

func (p *PingHandler) Handle(ctx *gin.Context) {
	p.logger.DebugLog("Server was just pinged")
	ctx.JSON(200, "pong")
}
