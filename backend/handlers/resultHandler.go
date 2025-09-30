package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/gin-gonic/gin"
)

type ResultHandler struct{ RouteHandlerBase }

func (r *ResultHandler) Init(logger logger.ILogger, modelHelper models.ISplitModelHelper, sessionManger managers.ISessionManager) {
	r.RouteHandlerBase.BaseInit(logger, modelHelper, sessionManger)

	r.pattern = `/ping`
	r.method = "GET"
	r.group = ""
	r.handler = r.Handle
}

func (r *ResultHandler) Handle(ctx *gin.Context) {

}

func testResult() IRouteHandlerBase {
	return &ResultHandler{}
}
