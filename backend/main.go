package main

import (
	"fmt"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/controllers/handlers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
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
		// concrete types
		utils.CreateJustProvider(utils.CreateGinEngine),
		utils.CreateJustProvider(utils.NewConnUpgrader),

		// interface types
		utils.CreateAnnotatedProvider(utils.NewLogger, new(logger.ILogger)),
		utils.CreateAnnotatedProvider(utils.NewDatabase, new(common.IDatabase)),
		utils.CreateAnnotatedProvider(utils.NewModelHelper, new(common.ISplitModelHelper)),
		utils.CreateAnnotatedProvider(utils.NewSessionManager, new(common.ISessionManager)),

		// Handlers
		utils.CreateHandlerProvider(utils.NewHandler[handlers.PingHandler]),
		utils.CreateHandlerProvider(utils.NewHandler[handlers.HomeHandler]),
		utils.CreateHandlerProvider(utils.NewHandler[handlers.SplitActionWsHandler]),
		utils.CreateHandlerProvider(utils.NewHandler[handlers.CloseHandler]),

		// route creator
		utils.CreateHandlerConsumer(utils.NewRouteCreater, new(routes.IRouteCreator)),

		// start the app
		fx.Invoke(utils.StartApp),
	).Run()
}
