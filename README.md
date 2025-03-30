# Test Case Generator API

This API service uses the Gemini API to generate test scenarios from user stories or requirements.

## Setup

1. Make sure you have Go installed on your system
2. Clone this repository
3. Install dependencies:

```bash
go mod tidy
```

## Running the Server

You can run the server in two ways:

1. Using environment variable:

```bash
export GEMINI_API_KEY=your_api_key_here
go run cmd/main.go
```

2. Using command line flag:

```bash
go run cmd/main.go -api-key=your_api_key_here
```

By default, the server runs on port 8080. You can change this using the `-port` flag:

```bash
go run cmd/main.go -port=3000 -api-key=your_api_key_here
```

## API Endpoints

### POST /generate-test-cases

Generates test scenarios from the provided description.

**Request Body:**

```json
{
  "contents": [
    {
      "parts": [
        {
          "text": "Your requirement or user story description here"
        }
      ]
    }
  ]
}
```

**Response:**

```json
{
  "scenarios": [
    {
      "id": "SC-LOGIN-01",
      "description": "Valid username and password",
      "inputData": "Correct username and password combination",
      "expected": "Successful login (redirects to appropriate page)"
    }
    // ... more scenarios
  ]
}
```

## Example Usage

```bash
curl -X POST \
  http://localhost:8080/generate-test-cases \
  -H 'Content-Type: application/json' \
  -d '{
    "contents": [{
      "parts": [{
        "text": "I am QA tester, I want you to generate test scenarios from following done items from developer: - Implement Login functionality"
      }]
    }]
  }'
```
