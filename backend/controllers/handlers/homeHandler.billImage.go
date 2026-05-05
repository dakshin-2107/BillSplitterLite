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
