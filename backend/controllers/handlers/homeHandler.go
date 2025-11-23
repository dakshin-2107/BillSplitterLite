package handlers

import (
	"fmt"
	"math"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Purpose :
- receiving the bill image and processing it AI.
- use the response from the LLM, add the split in the database
- return a (session ID, user ID) as part of a cookie and redirect the user to the session created
*/

type HomeHandler struct{ common.RouteHandlerBase }

func (h *HomeHandler) Init() {
	h.Pattern = `/home`
	h.Method = "POST"
	h.Group = ""
	h.Handler = h.Handle
}

func (h *HomeHandler) Handle(ctx *gin.Context) {

	ctx.Request.ParseForm()
	place := ctx.PostForm("place")
	splitDate, _ := time.Parse("02/01/2006", ctx.PostForm("dateTime"))
	names := ctx.PostFormArray("names")
	image, imageErr := ctx.FormFile("image")

	if imageErr == nil {

		newSplit := h.ModelHelper.CreateNewSplit()
		newSplit.CreatedAt = common.CreateSplitDate(time.Now())
		newSplit.LastUpdated = common.CreateSplitDate(time.Now())
		h.ParseItemsFromImage(image, newSplit)
		//h.ParseItemsFromImageDummy(image, newSplit)

		h.CreateParticipantsMap(names, newSplit)
		newSplit.BillID = GenerateNewBillId(place, splitDate)
		err := h.ModelHelper.AddBillIfNotExists(newSplit)
		if err == nil {

			ctx.SetCookie("BillsessionID", newSplit.BillID, math.MaxInt64, "/", "", true, true) // TODO: Update the cookie age limit
			ctx.SetCookie("UserID", "admin_boi", math.MaxInt64, "/", "", true, true)
			ctx.JSON(200, gin.H{
				"success": true,
				"message": "Session has been created",
				"bill":    newSplit,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"success": false,
		"message": "Parsing was not successful. Session could not be created.",
	})
}

func (h *HomeHandler) CreateParticipantsMap(names []string, split *common.Split) {

	currPersonCount := 1
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
		item.Takers = make(map[string]int, 0)
	}
}

func GenerateNewBillId(place string, splitDate time.Time) string {
	return uuid.New().String()
}

func testHome() common.IRouteHandlerBase {
	return &HomeHandler{}
}
