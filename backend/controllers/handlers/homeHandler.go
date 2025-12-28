package handlers

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"strings"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Purpose :
- receiving the bill image and processing it using AI.
- use the response from the LLM, adds the split in the database
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

	billID, billIdErr := ctx.Cookie(common.BILL_SESSION_ID_STR)
	if split, err := h.ModelHelper.GetBill(billID); billIdErr == nil && err == nil && split != nil {

		ctx.JSON(http.StatusOK, common.SessionAlreadyExistsResponse(*split))
		return
	} else {

		ctx.Request.ParseForm()
		place := ctx.PostForm("place")
		splitDate, _ := time.Parse("02/01/2006", ctx.PostForm("dateTime"))
		names := ctx.PostFormArray("names")
		image, imageErr := ctx.FormFile("image")

		if imageErr == nil {

			newSplit := h.ModelHelper.CreateNewSplit()
			newSplit.CreatedAt = time.Now()
			newSplit.LastUpdated = time.Now()
			h.ParseItemsFromImage(image, newSplit)
			//h.ParseItemsFromImageDummy(image, newSplit)

			h.CreateParticipantsMap(names, newSplit)
			newSplit.Date = common.CreateSplitDate(splitDate)
			newSplit.Location = place
			newSplit.BillID = GenerateNewBillId(place, splitDate)
			err := h.ModelHelper.AddBillIfNotExists(newSplit)
			if err == nil {
				// TODO: Update the cookie age limit
				ctx.SetCookie(common.BILL_SESSION_ID_STR, newSplit.BillID, math.MaxInt64, "/", "localhost", true, true)
				// TODO: Update the cookie age limit and extract the user name
				ctx.SetCookie(common.USER_SESSION_ID_STR, "admin_boi", math.MaxInt64, "/", "localhost", true, true)
				ctx.SetCookie(common.HELLO_THERE_STR, "General Kenobi", math.MaxInt64, "/", "localhost", false, false)
				ctx.JSON(http.StatusOK, common.SessionCreationSuccessResponse(*newSplit))
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, common.SessionCreationFailureResponse())
}

func (h *HomeHandler) CreateParticipantsMap(names []string, split *common.Split) {
	uniqueNames := make(map[string]bool)
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			uniqueNames[trimmed] = true
		}
	}

	participantsMap := make(map[string]string)
	usedIds := make(map[string]bool)

	// Process names in provided order for deterministic IDs
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" || !uniqueNames[trimmed] {
			continue
		}

		id := h.generateUniqueId(trimmed, usedIds)
		participantsMap[id] = trimmed
		usedIds[id] = true
		delete(uniqueNames, trimmed) // mark as processed
	}

	split.Participants = participantsMap

	for _, item := range split.Items {
		item.Takers = make(map[string]int)
	}
}

func (h *HomeHandler) generateUniqueId(name string, usedIds map[string]bool) string {
	words := strings.Fields(strings.ToUpper(name))
	if len(words) == 0 {
		return "USR"
	}

	firstWord := words[0]
	id := firstWord

	if !usedIds[id] {
		return id
	}

	// Collision resolution: append letters from the second word
	if len(words) > 1 {
		secondWord := words[1]
		currentId := id
		for i := 0; i < len(secondWord); i++ {
			currentId += string(secondWord[i])
			if !usedIds[currentId] {
				return currentId
			}
		}
		id = currentId // Use the full first word + full second word if still colliding
	}

	// Final fallback: append numbers if still not unique
	baseId := id
	for i := 1; i <= 99; i++ {
		candidate := fmt.Sprintf("%s%d", baseId, i)
		if !usedIds[candidate] {
			return candidate
		}
	}

	return uuid.New().String()[:5]
}

func GenerateNewBillId(place string, splitDate time.Time) string {
	return uuid.New().String()
}

func testHome() common.IRouteHandlerBase {
	return &HomeHandler{}
}
