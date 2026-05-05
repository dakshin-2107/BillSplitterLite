package handlers

import (
	"encoding/json"
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
Purpose:
- Developer-only endpoint (gated by EXPOSE_EXAMPLE_ENDPOINT env flag) that creates a
  split session from SAMPLE_JSON_STRING2 with a fixed set of dummy participants.
*/

const SAMPLE_JSON_STRING2 string = `{
  "bills": {
    "1": {
      "billId": 1,
      "total": 2592.00,
      "itemIdCounter": 11,
      "items": {
        "1": { "id": 1, "name": "Mango ale pint",    "price": 230.00, "takers": {} },
        "2": { "id": 2, "name": "Guava pint",         "price": 230.00, "takers": {} },
        "3": { "id": 3, "name": "Diet coke",          "price": 170.00, "takers": {} },
        "4": { "id": 4, "name": "Paneer ghee roast",  "price": 360.00, "takers": {} },
        "5": { "id": 5, "name": "Chicken pepper",     "price": 360.00, "takers": {} },
        "6": { "id": 6, "name": "Egg fried rice",     "price": 260.00, "takers": {} },
        "7": { "id": 7, "name": "Jeera rice",         "price": 190.00, "takers": {} },
        "8": { "id": 8, "name": "Tawa pulao pan",     "price": 340.00, "takers": {} },
        "9": { "id": 9, "name": "Murgh hariyali",     "price": 350.00, "takers": {} },
        "10": { "id": 10, "name": "Bill tax",         "price": 102.00, "takers": {} }
      }
    },
    "2": {
      "billId": 2,
      "total": 1350.00,
      "itemIdCounter": 6,
      "items": {
        "1": { "id": 1, "name": "Mango ale pint",    "price": 230.00, "takers": {} },
        "2": { "id": 2, "name": "Guava pint",         "price": 230.00, "takers": {} },
        "3": { "id": 3, "name": "Diet coke",          "price": 170.00, "takers": {} },
        "4": { "id": 4, "name": "Paneer ghee roast",  "price": 360.00, "takers": {} },
        "5": { "id": 5, "name": "Chicken pepper",     "price": 360.00, "takers": {} }
      }
    },
    "3": {
      "billId": 3,
      "total": 1800.00,
      "itemIdCounter": 8,
      "items": {
        "1": { "id": 1, "name": "Mango ale pint",    "price": 230.00, "takers": {} },
        "2": { "id": 2, "name": "Guava pint",         "price": 230.00, "takers": {} },
        "3": { "id": 3, "name": "Diet coke",          "price": 170.00, "takers": {} },
        "4": { "id": 4, "name": "Paneer ghee roast",  "price": 360.00, "takers": {} },
        "5": { "id": 5, "name": "Chicken pepper",     "price": 360.00, "takers": {} },
        "6": { "id": 6, "name": "Egg fried rice",     "price": 260.00, "takers": {} },
        "7": { "id": 7, "name": "Jeera rice",         "price": 190.00, "takers": {} }
      }
    }
  }
}`

type ExampleHandler struct{ common.RouteHandlerBase }

func (e *ExampleHandler) Init() {
	e.Pattern = `/example`
	e.Method = "GET"
	e.Group = ""
	e.Handler = e.Handle
}

func (e *ExampleHandler) Handle(ctx *gin.Context) {
	if os.Getenv("EXPOSE_EXAMPLE_ENDPOINT") != "true" {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Not found"})
		return
	}

	splitID, splitIdErr := ctx.Cookie(common.BILL_SESSION_ID_STR)
	if splitIdErr == nil {
		if split, err := e.ModelHelper.GetSplit(splitID); err == nil && split != nil {
			ctx.JSON(http.StatusOK, common.SessionAlreadyExistsResponse(*split))
			return
		}
	}

	newSplit := e.ModelHelper.CreateNewSplit()

	parsedBillSplit := strings.TrimPrefix(SAMPLE_JSON_STRING2, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")
	if err := json.Unmarshal([]byte(parsedBillSplit), newSplit); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to parse example data"})
		return
	}

	dummyNames := []string{"Kevin", "Dwight", "Pam", "Michael", "Jim", "Kelly", "Angela"}
	dummyLocNames := []string{"Poor Richard's", "Chilli's", "Hooters", "Schrute Farms", "Alfredo's pizza"}

	participants := make(map[string]string)
	for _, name := range dummyNames {
		participants[strings.ToUpper(name)] = name
	}
	newSplit.Participants = participants

	newSplit.SplitID = uuid.New().String()
	newSplit.CreatedAt = time.Now()
	newSplit.LastUpdated = time.Now()
	newSplit.BillIdCounter = len(newSplit.Bills) + 1

	locIndex := 0
	for billId, bill := range newSplit.Bills {
		bill.Location = dummyLocNames[locIndex]
		bill.Date = time.Now().Format("2006-01-02")
		newSplit.Bills[billId] = bill
		locIndex++
	}

	if err := e.ModelHelper.AddSplitIfNotExists(newSplit); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create example session"})
		return
	}

	domain := os.Getenv("BASE_URL")
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(common.BILL_SESSION_ID_STR, newSplit.SplitID, math.MaxInt64, "/", domain, false, true)
	ctx.SetCookie(common.USER_SESSION_ID_STR, "admin_boi", math.MaxInt64, "/", domain, false, true)
	ctx.JSON(http.StatusOK, common.SessionCreationSuccessResponse(*newSplit))
}
