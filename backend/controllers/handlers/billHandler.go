package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

/*
Purpose :
- Simple gets a bill using bill/session ID
*/

type BillHandler struct{ common.RouteHandlerBase }

func (b *BillHandler) Init() {
	b.Pattern = `/bill`
	b.Method = "GET"
	b.Group = ""
	b.Handler = b.Handle
}

func (b *BillHandler) Handle(ctx *gin.Context) {

	billId, err := ctx.Cookie("BillsessionID")
	if err != nil {
		ctx.JSON(200, gin.H{
			"success": false,
			"message": "Missing bill session ID cookie.",
		})
		return
	}

	ctx.Request.ParseForm()
	split, err := b.ModelHelper.GetBill(billId)

	if err != nil {
		ctx.JSON(200, gin.H{
			"success": false,
			"message": "Could not find bill requested ID.",
		})
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": split,
	})
}

func testSplitEditorHandler() common.IRouteHandlerBase {
	return &BillHandler{}
}
