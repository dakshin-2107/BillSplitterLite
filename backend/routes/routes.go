package routes

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/handlers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"

	"github.com/gin-gonic/gin"
)

func test() IRouteCreator {
	return &RouteCreator{}
}

type IRouteCreator interface {
	GetHandlerCount() int
	CreateRoute(handlersHolder gin.IRoutes, handler handlers.IRouteHandlerBase)
	Init(logger logger.ILogger, r *gin.Engine, modelHelper models.ISplitModelHelper, sessionManager managers.ISessionManager, handlers ...handlers.IRouteHandlerBase)
}

type RouteCreator struct {
	logger          logger.ILogger
	ginEngine       *gin.Engine
	modelHelper     models.ISplitModifier
	HandlerCount    int
	routerGroupDict map[string]*gin.RouterGroup
}

func (routeCreator *RouteCreator) Init(logger logger.ILogger, r *gin.Engine, modelHelper models.ISplitModelHelper, sessionManager managers.ISessionManager, handlers ...handlers.IRouteHandlerBase) {

	routeCreator.logger = logger
	routeCreator.ginEngine = r
	routeCreator.modelHelper = modelHelper
	routeCreator.routerGroupDict = make(map[string]*gin.RouterGroup)

	// middleware handlers
	for _, handler := range handlers {
		handler.Init(logger, modelHelper, sessionManager)
		if handler.Method() == "" {
			logger.DebugLog("Creating middleware group with pattern: " + handler.Pattern())
			routeCreator.routerGroupDict[handler.Pattern()] = r.Group(handler.Pattern(), handler.Handle)
		}
	}

	// non middleware handlers
	for _, handler := range handlers {
		if handler.Method() != "" {
			if handler.Group() != "" {
				routeHolder, exists := routeCreator.routerGroupDict[handler.Group()]
				if exists {
					logger.DebugLog("Initializing handler " + handler.Pattern() + " with method " + handler.Method() + " in group " + handler.Group())
					routeCreator.CreateRoute(routeHolder, handler)
				}
			} else {
				logger.DebugLog("Initializing handler " + handler.Pattern() + " with method " + handler.Method())
				routeCreator.CreateRoute(routeCreator.ginEngine, handler)
			}
		}

		routeCreator.HandlerCount++
	}
}

func (routeCreator *RouteCreator) GetHandlerCount() int {
	return routeCreator.HandlerCount
}

func (routeCreator *RouteCreator) CreateRoute(routeHolder gin.IRoutes, routeBase handlers.IRouteHandlerBase) {
	switch routeBase.Method() {
	case "GET":
		routeHolder.GET(routeBase.Pattern(), routeBase.Handle)
	case "POST":
		routeHolder.POST(routeBase.Pattern(), routeBase.Handle)
	case "PUT":
		routeHolder.PUT(routeBase.Pattern(), routeBase.Handle)
	case "DELETE":
		routeHolder.DELETE(routeBase.Pattern(), routeBase.Handle)
	}
}
