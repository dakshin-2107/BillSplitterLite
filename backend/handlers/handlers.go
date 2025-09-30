package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/gin-gonic/gin"
)

type IRouteHandlerBase interface {
	Group() string
	Method() string
	Pattern() string
	Handle(context *gin.Context)
	Init(logger logger.ILogger, modelHelper models.ISplitModelHelper, sessionManger managers.ISessionManager)
}

type RouteHandlerBase struct {
	group          string
	method         string
	pattern        string
	logger         logger.ILogger
	modelHelper    models.ISplitModelHelper
	handler        gin.HandlerFunc
	sessionManager managers.ISessionManager
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

func (routeBase *RouteHandlerBase) BaseInit(logger logger.ILogger, modelHelper models.ISplitModelHelper, sessionManger managers.ISessionManager) {
	routeBase.logger = logger
	routeBase.modelHelper = modelHelper
	routeBase.sessionManager = sessionManger
}
