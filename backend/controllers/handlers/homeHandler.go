package handlers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Purpose :
- receiving the bill image and processing it using AI
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

	splitID, splitIdErr := ctx.Cookie(common.BILL_SESSION_ID_STR)
	if splitIdErr == nil {
		if split, err := h.ModelHelper.GetSplit(splitID); err == nil && split != nil {
			ctx.JSON(http.StatusOK, common.SessionAlreadyExistsResponse(*split))
			return
		}
	}

	form, formErr := ctx.MultipartForm()
	if formErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to parse form"})
		return
	}

	locations := form.Value["locations"]
	dates := form.Value["dates"]
	names := form.Value["people"]
	images := form.File["images"]

	if len(images) > 0 && len(images) == len(locations) && len(images) == len(dates) && len(names) > 0 {

		newSplit := h.ModelHelper.CreateNewSplit()

		//err := h.ParseItemsFromImage(images, newSplit)
		err := h.ParseItemsFromImageDummy(images, newSplit)

		if err != nil || len(newSplit.Bills) != len(images) {
			h.Logger.InfoLog(fmt.Sprintf("Failed to parse bill with err: %v", err))
			return
		} else {

			h.CreateParticipantsMap(names, newSplit)
			newSplit.SplitID = uuid.New().String()
			newSplit.CreatedAt = time.Now()
			newSplit.LastUpdated = time.Now()
			newSplit.BillIdCounter = len(newSplit.Bills) + 1

			for billId, bill := range newSplit.Bills {
				bill.Location = locations[billId-1]
				bill.Date = dates[billId-1]
				newSplit.Bills[billId] = bill
				newSplit.TotalAmount += bill.TotalAmount
			}

			// add split to DB
			err = h.ModelHelper.AddSplitIfNotExists(newSplit)
			if err == nil {
				domain := os.Getenv("BASE_URL")
				ctx.SetSameSite(http.SameSiteStrictMode)
				ctx.SetCookie(common.BILL_SESSION_ID_STR, newSplit.SplitID, math.MaxInt64, "/", domain, false, true)
				ctx.SetCookie(common.USER_SESSION_ID_STR, "admin_boi", math.MaxInt64, "/", domain, false, true)
				ctx.JSON(http.StatusOK, common.SessionCreationSuccessResponse(*newSplit))
				return
			}
		}
	}

	domain := os.Getenv("BASE_URL")
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(common.BILL_SESSION_ID_STR, "", -1, "/", domain, false, true)
	ctx.SetCookie(common.USER_SESSION_ID_STR, "", -1, "/", domain, false, true)
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

func testHome() common.IRouteHandlerBase {
	return &HomeHandler{}
}
