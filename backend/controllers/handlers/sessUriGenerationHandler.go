package handlers

import (
	"net/http"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
)

type SessUriGenerationHandler struct {
	common.SessionRouteHandlerBase
}

func (s *SessUriGenerationHandler) Init() {
	s.Pattern = `/generate`
	s.Method = "GET"
	s.Group = ""
	s.Handler = s.Handle
}

func (s *SessUriGenerationHandler) Handle(ctx *gin.Context) {
	billID, err := ctx.Cookie(common.BILL_SESSION_ID_STR)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.MissingCookieResponse())
		return
	}

	ctx.JSON(http.StatusOK, common.SessionUriResponse(billID))
}
