package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/controllers/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/models/database"
	"github.com/dakshin-2107/BillSplitterLite/backend/models/modelHelper"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/fx"
)

func getAllowedOrigins() []string {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

func CreateGinEngine(logger logger.ILogger) *gin.Engine {
	logger.DebugLog("Instantiated gin engine")
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(),
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

func NewRouteCreater(logger logger.ILogger, modelHelper common.ISplitModelHelper, r *gin.Engine, sessionManager common.ISessionManager, routesList ...common.IRouteHandlerBase) IRouteCreator {
	rc := new(RouteCreator)
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
			return slices.Contains(getAllowedOrigins(), r.Header.Get("Origin"))
		},
	}
}

func StartApp(lc fx.Lifecycle, r *gin.Engine, logger logger.ILogger, rc IRouteCreator, db common.IDatabase, sessionManager common.ISessionManager) {
	logger.DebugLog("Starting the app")

	if !db.IsConnected() {
		panic("could not connect to redis")
	}

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	server := &http.Server{Addr: addr, Handler: r}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoLog(fmt.Sprintf("Server listening on %s", addr))
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.InfoLog(fmt.Sprintf("Server error: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoLog("Shutting down server...")
			if err := server.Shutdown(ctx); err != nil {
				logger.InfoLog(fmt.Sprintf("Server shutdown error: %v", err))
			}
			sessionManager.CleanSessions()
			db.DisconnectFromDatabase()
			logger.InfoLog("Shutdown complete.")
			return nil
		},
	})
}
