package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gemini-ai/pkg/models"
)

const (
	GeminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-8b:generateContent"
)

type Handler struct {
	apiKey string
}

func NewHandler(apiKey string) *Handler {
	return &Handler{
		apiKey: apiKey,
	}
}

func (h *Handler) GenerateTestCases(w http.ResponseWriter, r *http.Request) {
	// Create default request for test case generation
	req := models.GeminiRequest{
		Contents: []models.Content{
			{
				Parts: []models.Part{
					{
						Text: "I am QA tester, I want you to generate test scenarios from following done items from developer: - Implement Login functionality",
					},
				},
			},
		},
	}

	// Call Gemini API
	geminiResp, err := h.callGeminiAPI(req)
	if err != nil {
		http.Error(w, "Error calling Gemini API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the response and extract test scenarios
	scenarios := h.parseTestScenarios(geminiResp)

	// Send response
	response := models.TestScenarioResponse{
		Scenarios: scenarios,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) callGeminiAPI(req models.GeminiRequest) (*models.GeminiResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", GeminiAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-goog-api-key", h.apiKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, err
	}

	return &geminiResp, nil
}

func (h *Handler) parseTestScenarios(resp *models.GeminiResponse) []models.TestScenario {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return []models.TestScenario{}
	}

	text := resp.Candidates[0].Content.Parts[0].Text
	var scenarios []models.TestScenario

	// Extract scenarios from the text that are in table format
	lines := bytes.Split([]byte(text), []byte("\n"))
	for i := 0; i < len(lines); i++ {
		line := string(lines[i])
		if len(line) == 0 {
			continue
		}

		// Look for lines that start with "SC-LOGIN-" or similar pattern
		if len(line) > 2 && line[0] == '|' {
			parts := bytes.Split([]byte(line), []byte("|"))
			if len(parts) >= 5 {
				scenario := models.TestScenario{
					ID:          cleanString(parts[1]),
					Description: cleanString(parts[2]),
					InputData:   cleanString(parts[3]),
					Expected:    cleanString(parts[4]),
				}
				if scenario.ID != "Scenario ID" && scenario.ID != "" { // Skip header row
					scenarios = append(scenarios, scenario)
				}
			}
		}
	}

	return scenarios
}

func cleanString(s []byte) string {
	return string(bytes.TrimSpace(s))
}
