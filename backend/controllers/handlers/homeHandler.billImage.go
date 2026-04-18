package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"google.golang.org/genai"
)

/*
 Purpose :
 - Consists of helper methods that call the LLM and process the bill image
*/

// use gemini to parse the images
func (h *HomeHandler) ParseItemsFromImage(images []*multipart.FileHeader, split *common.Split) error {

	response, err := h.CallGeminiAPI(images)
	if err != nil {
		return err
	}

	parsedBillSplit := response.Candidates[0].Content.Parts[0].Text
	parsedBillSplit = strings.TrimPrefix(parsedBillSplit, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")

	h.Logger.DebugLog(parsedBillSplit)
	err = json.Unmarshal([]byte(parsedBillSplit), split)
	if err != nil {
		return err
	}

	return nil
}

func (h *HomeHandler) ParseItemsFromImageDummy(images []*multipart.FileHeader, split *common.Split) error {

	parsedBillSplit := SAMPLE_JSON_STRING2
	parsedBillSplit = strings.TrimPrefix(parsedBillSplit, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")

	h.Logger.DebugLog(parsedBillSplit)
	err := json.Unmarshal([]byte(parsedBillSplit), split)
	if err != nil {
		return err
	}

	return nil
}

func (h *HomeHandler) CallGeminiAPI(images []*multipart.FileHeader) (*genai.GenerateContentResponse, error) {
	var thinking_cost int32 = 0
	ctx := context.Background()

	// auto uses the API key from the env variable `GEMINI_API_KEY` automatically.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	parts := []*genai.Part{}
	for _, image := range images {
		file, err := image.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open image: %w", err)
		}
		defer file.Close()

		imageBytes := make([]byte, image.Size)
		_, err = file.Read(imageBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read image bytes: %w", err)
		}

		parts = append(parts, []*genai.Part{
			genai.NewPartFromBytes(imageBytes, "image/jpeg"),
			genai.NewPartFromText(os.Getenv("AI_PROMPT")),
		}...)

	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinking_cost,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Gemini GenerateContent failed: %w", err)
	}

	return result, nil
}

const SAMPLE_JSON_STRING2 string = `{
  "bills": {
    "1": {
      "billId": 1,
      "total": 2592.00,
      "itemIdCounter": 11,
      "items": {
        "1": {
          "id": 1,
          "name": "Mango ale pint",
          "price": 230.00,
          "takers": {}
        },
        "2": {
          "id": 2,
          "name": "Guava pint",
          "price": 230.00,
          "takers": {}
        },
        "3": {
          "id": 3,
          "name": "Diet coke",
          "price": 170.00,
          "takers": {}
        },
        "4": {
          "id": 4,
          "name": "Paneer ghee roast",
          "price": 360.00,
          "takers": {}
        },
        "5": {
          "id": 5,
          "name": "Chicken pepper",
          "price": 360.00,
          "takers": {}
        },
        "6": {
          "id": 6,
          "name": "Egg fried rice",
          "price": 260.00,
          "takers": {}
        },
        "7": {
          "id": 7,
          "name": "Jeera rice",
          "price": 190.00,
          "takers": {}
        },
        "8": {
          "id": 8,
          "name": "Tawa pulao pan",
          "price": 340.00,
          "takers": {}
        },
        "9": {
          "id": 9,
          "name": "Murgh hariyali",
          "price": 350.00,
          "takers": {}
        },
        "10": {
          "id": 10,
          "name": "Bill tax",
          "price": 102.00,
          "takers": {}
        }
      }
    }
  }
}`

