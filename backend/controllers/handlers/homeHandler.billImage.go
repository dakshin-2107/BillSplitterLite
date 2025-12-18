package handlers

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"strings"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"google.golang.org/genai"
)

/*
 Purpose :
 - Consists of helper methods that call the LLM and process the bill image
*/

const AI_Prompt string = `Using the image provided of a bill and return the following details from the bill in the below JSON format. For location get the only name of the place, do not get the entire address. If you cannot find the whole name just provide an empty string.
{
	"date": <date on the bill>,
	"location": <location on the bill>,
	"total": <total amount>
	"items": {
        "id" : {
            "id": <id in increasing order as a string>,
            "name": <name of the item>
            "price": <price of the item>,
            "takers": <always an empty dictionary>
        },
    }
}
`

const SAMPLE_JSON_STRING string = `{
  "date": "04-Oct-2025",
  "location": "Namma biryani",
  "total": 250.00,
  "items": {
    "1": {
      "id": "1",
      "name": "INDIAN GRILL CKN HALF",
      "price": 200.00,
      "takers":  {}
    },
    "2": {
      "id": "2",
      "name": "BIRIYANI RICE HALF",
      "price": 50.00,
      "takers": {}
    }
  }
}`

// use gemini to parse the images
func (h *HomeHandler) ParseItemsFromImage(image *multipart.FileHeader, billSplit *common.Split) error {

	response, err := h.CallGeminiAPI(image)
	if err != nil {
		return err
	}

	parsedBillSplit := response.Candidates[0].Content.Parts[0].Text
	parsedBillSplit = strings.TrimPrefix(parsedBillSplit, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")

	h.Logger.DebugLog(parsedBillSplit)
	err = json.Unmarshal([]byte(parsedBillSplit), billSplit)
	if err != nil {
		return err
	}

	return nil
}

func (h *HomeHandler) ParseItemsFromImageDummy(image *multipart.FileHeader, billSplit *common.Split) error {

	parsedBillSplit := SAMPLE_JSON_STRING
	parsedBillSplit = strings.TrimPrefix(parsedBillSplit, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")

	h.Logger.DebugLog(parsedBillSplit)
	err := json.Unmarshal([]byte(parsedBillSplit), billSplit)
	if err != nil {
		return err
	}

	return nil
}

func (h *HomeHandler) CallGeminiAPI(image *multipart.FileHeader) (*genai.GenerateContentResponse, error) {
	var thinking_cost int32 = 0
	ctx := context.Background()

	// auto uses the API key from the env variable `GEMINI_API_KEY` automatically.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	file, err := image.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	imageBytes := make([]byte, image.Size)
	_, err = file.Read(imageBytes)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/jpeg"),
		genai.NewPartFromText(AI_Prompt),
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
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}
