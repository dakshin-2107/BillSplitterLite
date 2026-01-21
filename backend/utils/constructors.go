package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/controllers/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/models/database"
	"github.com/dakshin-2107/BillSplitterLite/backend/models/modelHelper"
	"github.com/dakshin-2107/BillSplitterLite/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var allowedOrigins = []string{
	"http://localhost:5173",
	"http://localhost:5174",
	"http://192.168.0.169:5173",
	"http://192.168.0.169:5174",
	"https://splitzo.dak-shin.com",
}

func CreateGinEngine(logger logger.ILogger) *gin.Engine {
	logger.DebugLog("Instantiated gin engine")
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return r
}

func NewLogger() logger.ILogger {
	logger := new(logger.Logger)
	logger.Init()
	logger.DebugLog("Instantiated logger")
	return logger
}

func NewDatabase(logger logger.ILogger) (common.IDatabase, error) {
	redisDB := new(database.RedisDatabaseConnection)
	err := redisDB.Init(logger)
	logger.DebugLog("Instantiated database connection")
	return redisDB, err
}

func NewModelHelper(db common.IDatabase, logger logger.ILogger) (common.ISplitModelHelper, error) {
	helper := new(modelHelper.SplitModelHelper)
	helper.Init(db, logger)
	logger.DebugLog("Instantiated model helper")
	return helper, nil
}

func NewHandler[T any]() common.IRouteHandlerBase {

	// So that I dont waste time again.
	// The below code does not work because Go does not allow direct type assertions on generic types as shown below.
	// handler, ok := new(T).(IRouteHandlerBase)

	handler, ok := any(new(T)).(common.IRouteHandlerBase)
	if !ok {
		panic("T must implement IRouteHandlerBase")
	}

	return handler
}

func NewRouteCreater(logger logger.ILogger, modelHelper common.ISplitModelHelper, r *gin.Engine, sessionManager common.ISessionManager, routesList ...common.IRouteHandlerBase) routes.IRouteCreator {
	rc := new(routes.RouteCreator)
	logger.DebugLog("Instantiated route creator")

	rc.Init(logger, r, modelHelper, sessionManager, routesList...)
	logger.DebugLog(fmt.Sprintf("Route creator created %d routes", rc.GetHandlerCount()))
	return rc
}

func NewSessionManager(logger logger.ILogger, connUpgrader *websocket.Upgrader, modelHelper common.ISplitModelHelper) common.ISessionManager {
	sessionMananger := new(managers.SessionManager)
	sessionMananger.Init(logger, connUpgrader, modelHelper)
	return sessionMananger
}

func NewConnUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
	}
}

func StartApp(r *gin.Engine, logger logger.ILogger, rc routes.IRouteCreator, db common.IDatabase) {
	logger.DebugLog("Starting the app")

	if !db.IsConnected() {
		panic("could not connect to redis")
	}

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
