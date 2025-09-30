package main

import (
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/handlers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/dakshin-2107/BillSplitterLite/backend/routes"
	"github.com/dakshin-2107/BillSplitterLite/backend/utils"
	"github.com/joho/godotenv"

	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load any environment files")
		//panic(err)
	}

	fx.New(
		utils.CreateJustProvider(utils.CreateGinEngine),
		utils.CreateJustProvider(utils.NewConnUpgrader),
		utils.CreateAnnotatedProvider(utils.NewLogger, new(logger.ILogger)),
		utils.CreateAnnotatedProvider(utils.NewDatabase, new(models.IDatabase)),
		utils.CreateAnnotatedProvider(utils.NewModelHelper, new(models.ISplitModelHelper)),
		utils.CreateAnnotatedProvider(utils.NewSessionManager, new(managers.ISessionManager)),

		// Handlers
		utils.CreateHandlerProvider(utils.NewHandler[handlers.PingHandler]),
		utils.CreateHandlerProvider(utils.NewHandler[handlers.HomeHandler]),

		// route creator
		utils.CreateHandlerConsumer(utils.NewRouteCreater, new(routes.IRouteCreator)),

		// start the app
		fx.Invoke(utils.StartApp),
	).Run()
}
