package handlers

import (
	"fmt"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Purpose :
- receiving the bill image and processing it AI.
- use the response from the MCP agent, add the split in the database
- return a (session ID, user ID) as part of a cookie and redirect the user to the session created
*/

type HomeHandler struct{ RouteHandlerBase }

func (h *HomeHandler) Init(logger logger.ILogger, modelHelper models.ISplitModelHelper, sessionManger managers.ISessionManager) {
	h.RouteHandlerBase.BaseInit(logger, modelHelper, sessionManger)

	h.pattern = `/home`
	h.method = "POST"
	h.group = ""
	h.handler = h.Handle
}

func (h *HomeHandler) Handle(ctx *gin.Context) {

	ctx.Request.ParseForm()
	place := ctx.PostForm("place")
	splitDate, _ := time.Parse("02/01/2006", ctx.PostForm("dateTime"))
	names := ctx.PostFormArray("names")
	image, imageErr := ctx.FormFile("image")

	if imageErr == nil {

		newSplit := h.modelHelper.CreateNewSplit()
		newSplit.CreatedAt = h.modelHelper.CreateSplitDate(time.Now())
		newSplit.LastUpdated = h.modelHelper.CreateSplitDate(time.Now())
		h.ParseItemsFromImage(image, newSplit)
		//h.ParseItemsFromImageDummy(image, newSplit)

		h.CreateParticipantsMap(names, newSplit)
		newSplit.BillID = GenerateNewBillId(place, splitDate)
		err := h.modelHelper.AddBillIfNotExists(newSplit)
		if err == nil {

			ctx.SetCookie("BillsessionID", newSplit.BillID, 3600, "/", "", true, true)
			ctx.JSON(200, gin.H{
				"success": true,
				"message": "Session has been created",
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"success": false,
		"message": "Parsing was not successful. Session could not be created.",
	})
}

func (h *HomeHandler) CreateParticipantsMap(names []string, split *models.Split) {

	currPersonCount := 1

	// ensure only unique names
	filteredNames := make(map[string]string)
	for _, name := range names {
		filteredNames[name] = ""
	}

	participantsMap := make(map[string]string)
	for name, _ := range filteredNames {
		participantsMap[fmt.Sprintf("%d", currPersonCount)] = name
		currPersonCount++
	}

	split.Participants = participantsMap

	for _, item := range split.Items {
		item.Takers = make([]string, 0)
	}
}

func GenerateNewBillId(place string, splitDate time.Time) string {
	return uuid.New().String()
}

func testHome() IRouteHandlerBase {
	return &HomeHandler{}
}
