package common

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gin-gonic/gin"
)

type RouteHandlerBase struct {
	Group       string // used for adding middle ware, method field must be empty if it's a middleware
	Method      string
	Pattern     string
	Logger      logger.ILogger
	ModelHelper ISplitModelHelper
	Handler     gin.HandlerFunc
}

type IRouteHandlerBase interface {
	GetGroup() string
	GetMethod() string
	GetPattern() string
	Handle(context *gin.Context)
	Init()
	BaseInit(logger logger.ILogger, modelHelper ISplitModelHelper)
}

func (routeBase *RouteHandlerBase) GetMethod() string {
	return routeBase.Method
}

func (routeBase *RouteHandlerBase) GetPattern() string {
	return routeBase.Pattern
}

func (routeBase *RouteHandlerBase) GetGroup() string {
	return routeBase.Group
}

func (routeBase *RouteHandlerBase) BaseInit(logger logger.ILogger, modelHelper ISplitModelHelper) {
	routeBase.Logger = logger
	routeBase.ModelHelper = modelHelper
}

// Session handler base - has session manager
type SessionRouteHandlerBase struct {
	RouteHandlerBase
	Logger         logger.ILogger
	ModelHelper    ISplitModelHelper
	SessionManager ISessionManager
}

type ISessionRouteHandlerBase interface {
	SessionHandlerBaseInit(logger logger.ILogger, modelHelper ISplitModelHelper, sessionManger ISessionManager)
}

func (routeSessionBase *SessionRouteHandlerBase) SessionHandlerBaseInit(logger logger.ILogger, modelHelper ISplitModelHelper, sessionManger ISessionManager) {
	routeSessionBase.SessionManager = sessionManger
	routeSessionBase.ModelHelper = modelHelper
	routeSessionBase.Logger = logger
}
