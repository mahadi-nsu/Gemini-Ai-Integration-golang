package models

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response structure from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

// TestScenario represents a single test scenario
type TestScenario struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	InputData   string `json:"inputData"`
	Expected    string `json:"expected"`
}

// TestScenarioResponse represents the formatted response we'll send to clients
type TestScenarioResponse struct {
	Scenarios []TestScenario `json:"scenarios"`
}
