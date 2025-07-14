package utils

import (
	"fmt"
	"os"

	"github.com/dakshin-2107/BillSplitterLite/backend/handlers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/dakshin-2107/BillSplitterLite/backend/routes"
	"github.com/gin-gonic/gin"
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
	mongo := new(models.RedisDatabaseConnection)
	err := mongo.Init(logger)
	logger.DebugLog("Instantiated database connection")
	return mongo, err
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

func NewRouteCreater(logger logger.ILogger, modelHelper models.ISplitModelHelper, r *gin.Engine, routesList ...handlers.IRouteHandlerBase) routes.IRouteCreator {
	rc := new(routes.RouteCreator)
	logger.DebugLog("Instantiated route creator")

	rc.Init(logger, r, modelHelper, routesList...)
	logger.DebugLog(fmt.Sprintf("Route creator created %d routes", rc.GetHandlerCount()))
	return rc
}

func StartApp(r *gin.Engine, logger logger.ILogger, rc routes.IRouteCreator) {
	logger.DebugLog("Starting the app")
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
