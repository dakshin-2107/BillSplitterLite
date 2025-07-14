package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/gin-gonic/gin"
)

type IRouteHandlerBase interface {
	Group() string
	Method() string
	Pattern() string
	Handle(context *gin.Context)
	Init(logger logger.ILogger, modelHelper models.ISplitModelHelper)
}

type RouteHandlerBase struct {
	group       string
	method      string
	pattern     string
	logger      logger.ILogger
	modelHelper models.ISplitModelHelper
	handler     gin.HandlerFunc
}

func (routeBase *RouteHandlerBase) Method() string {
	return routeBase.method
}

func (routeBase *RouteHandlerBase) Pattern() string {
	return routeBase.pattern
}

func (routeBase *RouteHandlerBase) Group() string {
	return routeBase.group
}
