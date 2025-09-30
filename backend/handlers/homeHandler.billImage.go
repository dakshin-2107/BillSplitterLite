package handlers

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"strings"

	"github.com/dakshin-2107/BillSplitterLite/backend/models"
	"google.golang.org/genai"
)

const AI_Prompt string = `Using the image provided of a bill and return the following details from the bill in the below JSON format. For location get the only name of the place, do not get the entire address. If you cannot find the whole name just provide an empty string.
{
	"date": <date on the bill>,
	"location": <location on the bill>,
	"total": <total amount>
	"items": [
		{
			"id": <id in increasing order as a string>,
			"name": <name of the item>
			"price": <price of the item>,
			"takers": <always an empty list>
		},
	]
}
`

const SAMPLE_JSON_STRING string = `{
    "date": "31/08/24",
    "location": "RR Nagar + High Drate",
    "total": 4037.00,
    "items": [
        {
            "id": "1",
            "name": "Guava Chilli",
            "price": 235.00,
            "takers": []
        },
        {
            "id": "2",
            "name": "Cheese Burst Nachos Paneer",
            "price": 360.00,
            "takers": []
        },
        {
            "id": "3",
            "name": "Sprite",
            "price": 240.00,
            "takers": []
        },
        {
            "id": "4",
            "name": "Diet Coke",
            "price": 240.00,
            "takers": []
        },
        {
            "id": "5",
            "name": "Scary House Mushroom Pepperdry",
            "price": 310.00,
            "takers": []
        },
        {
            "id": "6",
            "name": "Kalmi Kabab (2Pcs)",
            "price": 360.00,
            "takers": []
        },
        {
            "id": "7",
            "name": "Paneer Tikka",
            "price": 200.00,
            "takers": []
        },
        {
            "id": "8",
            "name": "Virgin Calada",
            "price": 235.00,
            "takers": []
        },
        {
            "id": "9",
            "name": "Sweet 16",
            "price": 470.00,
            "takers": []
        },
        {
            "id": "10",
            "name": "Farm House Pizza",
            "price": 195.00,
            "takers": []
        },
        {
            "id": "11",
            "name": "Mexican Pizza",
            "price": 195.00,
            "takers": []
        },
        {
            "id": "12",
            "name": "Veg Burnt Garlic Rice",
            "price": 225.00,
            "takers": []
        },
        {
            "id": "13",
            "name": "Curd Rice",
            "price": 120.00,
            "takers": []
        },
        {
            "id": "14",
            "name": "Dal Fry",
            "price": 170.00,
            "takers": []
        },
        {
            "id": "15",
            "name": "Coke",
            "price": 60.00,
            "takers": []
        }
    ]
}`

// use gemini to parse the images
func (h *HomeHandler) ParseItemsFromImage(image *multipart.FileHeader, billSplit *models.Split) error {

	response, err := h.CallGeminiAPI(image)
	if err != nil {
		return err
	}

	parsedBillSplit := response.Candidates[0].Content.Parts[0].Text
	parsedBillSplit = strings.TrimPrefix(parsedBillSplit, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")

	h.logger.DebugLog(parsedBillSplit)
	err = json.Unmarshal([]byte(parsedBillSplit), billSplit)
	if err != nil {
		return err
	}

	return nil
}

func (h *HomeHandler) ParseItemsFromImageDummy(image *multipart.FileHeader, billSplit *models.Split) error {

	parsedBillSplit := SAMPLE_JSON_STRING
	parsedBillSplit = strings.TrimPrefix(parsedBillSplit, "```json")
	parsedBillSplit = strings.TrimSuffix(parsedBillSplit, "```")

	h.logger.DebugLog(parsedBillSplit)
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
