package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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
	// Parse the request body to get feature details
	var featureRequest struct {
		FeatureName string `json:"featureName"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&featureRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a structured prompt for test case generation
	prompt := `As a QA Engineer, I need comprehensive test scenarios for the following feature:

Feature: ${featureName}
Description: ${description}

Please generate detailed test scenarios covering:
1. Positive test cases (Happy paths)
2. Negative test cases (Error handling)
3. Validation test cases
4. Edge cases
5. Security considerations
6. Performance aspects

Format the scenarios in a table with the following columns:
| Scenario ID | Description | Test Data/Steps | Expected Result |

Note: 
- Scenario IDs should follow the format: SC-{FEATURE}-{NUMBER}
- Include specific test data and validation criteria
- Consider different user roles if applicable
- Include boundary conditions
- Consider integration points with other features`

	// Replace placeholders with actual values
	prompt = strings.ReplaceAll(prompt, "${featureName}", featureRequest.FeatureName)
	prompt = strings.ReplaceAll(prompt, "${description}", featureRequest.Description)

	// Create Gemini request
	req := models.GeminiRequest{
		Contents: []models.Content{
			{
				Parts: []models.Part{
					{
						Text: prompt,
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
