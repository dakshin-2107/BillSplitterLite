package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dakshin-2107/BillSplitterLite/backend/handlers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/dakshin-2107/BillSplitterLite/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func CreateGinEngine(logger logger.ILogger) *gin.Engine {
	logger.DebugLog("Instantiated gin engine")
	return gin.New()
}

func NewLogger() logger.ILogger {
	logger := new(logger.Logger)
	logger.Init()
	logger.DebugLog("Instantiated logger")
	return logger
}

func NewDatabase(logger logger.ILogger) (models.IDatabase, error) {
	redisDB := new(models.RedisDatabaseConnection)
	err := redisDB.Init(logger)
	logger.DebugLog("Instantiated database connection")
	return redisDB, err
}

func NewModelHelper(db models.IDatabase, logger logger.ILogger) (models.ISplitModelHelper, error) {
	helper := new(models.SplitModelHelper)
	helper.Init(db, logger)
	logger.DebugLog("Instantiated model helper")
	return helper, nil
}

func NewHandler[T any]() handlers.IRouteHandlerBase {

	// So that I dont waste time again.
	// The below code does not work because Go does not allow direct type assertions on generic types as shown below.
	// handler, ok := new(T).(IRouteHandlerBase)

	handler, ok := any(new(T)).(handlers.IRouteHandlerBase)
	if !ok {
		panic("T must implement IRouteHandlerBase")
	}

	return handler
}

func NewRouteCreater(logger logger.ILogger, modelHelper models.ISplitModelHelper, r *gin.Engine, sessionManager managers.ISessionManager, routesList ...handlers.IRouteHandlerBase) routes.IRouteCreator {
	rc := new(routes.RouteCreator)
	logger.DebugLog("Instantiated route creator")

	rc.Init(logger, r, modelHelper, sessionManager, routesList...)
	logger.DebugLog(fmt.Sprintf("Route creator created %d routes", rc.GetHandlerCount()))
	return rc
}

func NewSessionManager(logger logger.ILogger, connUpgrader *managers.ConnectionUpgrader) managers.ISessionManager {
	sessionMananger := new(managers.SessionManager)
	sessionMananger.Init(logger, connUpgrader)
	return sessionMananger
}

func NewConnUpgrader() *managers.ConnectionUpgrader {
	return &managers.ConnectionUpgrader{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func StartApp(r *gin.Engine, logger logger.ILogger, rc routes.IRouteCreator, db models.IDatabase) {
	logger.DebugLog("Starting the app")

	if !db.IsConnected() {
		panic("could not connect to redis")
	}

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
