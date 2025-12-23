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
			newSplit.CreatedAt = common.CreateSplitDate(time.Now())
			newSplit.LastUpdated = common.CreateSplitDate(time.Now())
			//h.ParseItemsFromImage(image, newSplit)
			h.ParseItemsFromImageDummy(image, newSplit)

			h.CreateParticipantsMap(names, newSplit)
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
		if strings.TrimSpace(name) != "" {
			uniqueNames[name] = true
		}
	}

	participantsMap := make(map[string]string)
	usedIds := make(map[string]bool)

	// Since maps don't guarantee order, and we want deterministic IDs (e.g., first one gets the 3-char prefix),
	// it might be better to iterate the original 'names' slice but only process unique ones.
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
	cleanName := ""
	for _, r := range strings.ToUpper(name) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			cleanName += string(r)
		}
	}

	if cleanName == "" {
		cleanName = "USR"
	}

	// Try first 3 chars
	id := cleanName
	if len(id) > 3 {
		id = id[:3]
	}
	if !usedIds[id] {
		return id
	}

	// Try first 2 chars + digit (1-9)
	base2 := cleanName
	if len(base2) > 2 {
		base2 = base2[:2]
	}
	for i := 1; i <= 9; i++ {
		candidate := fmt.Sprintf("%s%d", base2, i)
		if !usedIds[candidate] {
			return candidate
		}
	}

	// Try first char + 2 digits (01-99)
	base1 := string(cleanName[0])
	for i := 1; i <= 99; i++ {
		candidate := fmt.Sprintf("%s%02d", base1, i)
		if !usedIds[candidate] {
			return candidate
		}
	}

	// Fallback to serial P01, P02...
	for i := 1; i <= 99; i++ {
		candidate := fmt.Sprintf("P%02d", i)
		if !usedIds[candidate] {
			return candidate
		}
	}

	return uuid.New().String()[:3]
}

func GenerateNewBillId(place string, splitDate time.Time) string {
	return uuid.New().String()
}

func testHome() common.IRouteHandlerBase {
	return &HomeHandler{}
}
