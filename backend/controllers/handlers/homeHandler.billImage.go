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

const AI_Prompt string = `Using the multiple images of the bills provided, return the list of items from the bills in the below JSON format. 
Find the tax amount as well, it is usually written as Tax or GST or SGST/CGST or service charge or service tax or VAT or any other similar term. Add the tax items 
as a single item called "Bill tax" for each bill. GST is the sum of SGST and CGST. Do not add all 3, add only GST or the sum of CGST or SGST. Add any roundoff amount to the tax item.
The bill id should start from 1 and increment by 1 for successive bills. Within each bill, the item id should start from 1 and increment by 1 for each successive item.
In some bills for each item there is rate, quantity and amount. In that case, use only the amount for the price of the item. 
Normalize all item names to sentence case.
{
  "totalAmount": <sum of total field of all the bills>,
  "billIdCounter": <number of bills> + 1,
  "bills": {
    "<billId>": {
      "billId": <billId as a number>,
      "total": <number with 2 decimal places>,
	  "itemIdCounter": <number of items in this bill> + 1,
      "items": {
        "<number1>": {
          "id": <number same as the key in the items map>,
          "name": "<item name>",
          "price": <number with 2 decimal places>,
          "takers": { }
        },
        "<number2>": {
          "id": <number same as the key in the items map>,
          "name": "<item name>",
          "price": <number with 2 decimal places>,
          "takers": { }
        }
      }
    },
    "<billId2>": {
		"billId": <billId2 as a number>,
		"total": <number with 2 decimal places>,
		"itemIdCounter": <number of items in this bill> + 1,
		"items": {
			"<number1>": {
				"id": <number same as the key in the items map>,
				"name": "<item name>",
				"price": <number with 2 decimal places>,
				"takers": { }
			}
		}
    }
  }
}
`

const SAMPLE_JSON_STRING string = `{
	"totalAmount": 450.75,
	"participants": {
		"Dakshin": "Dakshin",
		"Kevin": "Kevin",
		"Angela": "Angela"
	},
	"bills": {
		"1": {
			"billId": 1,
			"total": 250.00,
			"itemIdCounter": 3,
			"items": {
				"1": {
					"id": 1,
					"name": "Indian Grill Ckn Half",
					"price": 200.00,
					"takers": {}
				},
				"2": {
					"id": 2,
					"name": "Biriyani Rice Half",
					"price": 50.00,
					"takers": {}
				}
			}
		},
		"2": {
			"billId": 2,
			"total": 200.75,
			"itemIdCounter": 2,
			"items": {
				"1": {
					"id": 1,
					"name": "Iced Americano",
					"price": 200.75,
					"takers": {}
				}
			}
		}
	}
}`

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

	parsedBillSplit := SAMPLE_JSON_STRING
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
		log.Fatal(err)
	}

	parts := []*genai.Part{}
	for _, image := range images {
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

		parts = append(parts, []*genai.Part{
			genai.NewPartFromBytes(imageBytes, "image/jpeg"),
			genai.NewPartFromText(AI_Prompt),
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
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}
