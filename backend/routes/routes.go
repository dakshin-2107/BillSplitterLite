package routes

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"

	"github.com/gin-gonic/gin"
)

func test() IRouteCreator {
	return &RouteCreator{}
}

type IRouteCreator interface {
	GetHandlerCount() int
	CreateRoute(handlersHolder gin.IRoutes, handler common.IRouteHandlerBase)
	Init(logger logger.ILogger, r *gin.Engine, modelHelper common.ISplitModelHelper, sessionManager common.ISessionManager, handlers ...common.IRouteHandlerBase)
}

type RouteCreator struct {
	logger          logger.ILogger
	ginEngine       *gin.Engine
	modelHelper     common.ISplitModifier
	HandlerCount    int
	routerGroupDict map[string]*gin.RouterGroup
}

func (routeCreator *RouteCreator) Init(logger logger.ILogger, r *gin.Engine, modelHelper common.ISplitModelHelper, sessionManager common.ISessionManager, handlersList ...common.IRouteHandlerBase) {

	routeCreator.logger = logger
	routeCreator.ginEngine = r
	routeCreator.modelHelper = modelHelper
	routeCreator.routerGroupDict = make(map[string]*gin.RouterGroup)

	// middleware handlers
	for _, handler := range handlersList {

		if sessionHandler, ok := handler.(common.ISessionRouteHandlerBase); ok {
			sessionHandler.SessionHandlerBaseInit(logger, modelHelper, sessionManager)
		}

		handler.BaseInit(logger, modelHelper)
		handler.Init()
		if handler.GetMethod() == "" {
			logger.DebugLog("Creating middleware group with pattern: " + handler.GetPattern())
			routeCreator.routerGroupDict[handler.GetPattern()] = r.Group(handler.GetPattern(), handler.Handle)
		}
	}

	// non middleware handlers
	for _, handler := range handlersList {
		if handler.GetMethod() != "" {
			if handler.GetGroup() != "" {
				routeHolder, exists := routeCreator.routerGroupDict[handler.GetGroup()]
				if exists {
					logger.DebugLog("Initializing handler " + handler.GetPattern() + " with method " + handler.GetMethod() + " in group " + handler.GetGroup())
					routeCreator.CreateRoute(routeHolder, handler)
				}
			} else {
				logger.DebugLog("Initializing handler " + handler.GetPattern() + " with method " + handler.GetMethod())
				routeCreator.CreateRoute(routeCreator.ginEngine, handler)
			}
		}

		routeCreator.HandlerCount++
	}
}

func (routeCreator *RouteCreator) GetHandlerCount() int {
	return routeCreator.HandlerCount
}

func (routeCreator *RouteCreator) CreateRoute(routeHolder gin.IRoutes, routeBase common.IRouteHandlerBase) {
	switch routeBase.GetMethod() {
	case "GET":
		routeHolder.GET(routeBase.GetPattern(), routeBase.Handle)
	case "POST":
		routeHolder.POST(routeBase.GetPattern(), routeBase.Handle)
	case "PUT":
		routeHolder.PUT(routeBase.GetPattern(), routeBase.Handle)
	case "DELETE":
		routeHolder.DELETE(routeBase.GetPattern(), routeBase.Handle)
	}
}
