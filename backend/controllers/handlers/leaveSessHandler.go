package handlers

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

type LeaveSessionHandler struct {
	common.SessionRouteHandlerBase
}

func (ls *LeaveSessionHandler) Init() {
	ls.Pattern = `/leave`
	ls.Method = "GET"
	ls.Group = ""
	ls.Handler = ls.Handle
}

func (ls *LeaveSessionHandler) Handle(ctx *gin.Context) {

}
