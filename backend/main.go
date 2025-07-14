package main

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/handlers"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/dakshin-2107/BillSplitterLite/backend/routes"
	"github.com/dakshin-2107/BillSplitterLite/backend/utils"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fx.New(
		utils.CreateJustProvider(utils.CreateGinEngine),
		utils.CreateAnnotatedProvider(utils.NewLogger, new(logger.ILogger)),
		utils.CreateAnnotatedProvider(utils.NewDatabase, new(models.IDatabase)),
		utils.CreateAnnotatedProvider(utils.NewModelHelper, new(models.ISplitModelHelper)),

		// Handlers
		utils.CreateHandlerProvider(utils.NewHandler[handlers.PingHandler]),

		// route creator
		utils.CreateHandlerConsumer(utils.NewRouteCreater, new(routes.IRouteCreator)),

		// start the app
		fx.Invoke(utils.StartApp),
	).Run()
}
