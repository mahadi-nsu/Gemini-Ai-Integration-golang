# Test Case Generator API

This API service uses the Gemini API to generate test scenarios from user stories or requirements.

## Setup

1. Make sure you have Go installed on your system
2. Clone this repository
3. Copy `.env.example` to `.env` and add your Gemini API key
4. Install dependencies:

```bash
go mod tidy
```

## Running the Server

1. Using the run script:

```bash
./run.sh
```

The server runs on port 5001 by default.

## API Endpoints

### POST /generate-test-cases

Generates comprehensive test scenarios for a given feature.

**Request:**

- Method: `POST`
- URL: `http://localhost:5001/generate-test-cases`
- Headers:
  - `Content-Type: application/json`

**Request Body:**

```json
{
  "featureName": "User Registration",
  "description": "Allow new users to create an account with email, password, and profile information. The registration process includes email verification and password strength validation."
}
```

**Response Format:**

```json
{
  "scenarios": [
    {
      "id": "SC-REGISTER-001",
      "description": "Successful user registration with valid inputs",
      "inputData": "Email: valid@email.com, Password: StrongP@ss1, Name: John Doe",
      "expected": "User account created successfully, verification email sent"
    },
    {
      "id": "SC-REGISTER-002",
      "description": "Registration with invalid email format",
      "inputData": "Email: invalid.email, Password: StrongP@ss1",
      "expected": "Error message: 'Please enter a valid email address'"
    }
    // ... more scenarios
  ]
}
```

## Example Usage

### Basic Registration Feature

```json
{
  "featureName": "User Registration",
  "description": "Allow new users to create an account with following requirements:\n1. Email must be valid and unique\n2. Password must be at least 8 characters with 1 uppercase, 1 lowercase, 1 number, and 1 special character\n3. User must provide first name and last name\n4. User must agree to terms and conditions\n5. Email verification is required after registration"
}
```

### Login Feature

```json
{
  "featureName": "User Login",
  "description": "User login with following requirements:\n1. User can login with email and password\n2. System should lock account after 3 failed attempts\n3. Forgot password option should be available\n4. Remember me option should be available\n5. User should be redirected to dashboard after successful login"
}
```

### Password Reset Feature

```json
{
  "featureName": "Password Reset",
  "description": "Password reset functionality with following requirements:\n1. User can request password reset via email\n2. Reset link should expire after 1 hour\n3. New password must meet security requirements\n4. User should receive confirmation email after password change\n5. All active sessions should be terminated after password change"
}
```

## Testing with Postman

1. Create a new request in Postman
2. Set method to `POST`
3. Enter URL: `http://localhost:5001/generate-test-cases`
4. Add header:
   - Key: `Content-Type`
   - Value: `application/json`
5. In the Body tab:
   - Select `raw`
   - Select `JSON` from the dropdown
   - Paste one of the example request bodies above

## Note

The generated test scenarios will cover:

- Positive test cases (Happy paths)
- Negative test cases (Error handling)
- Validation test cases
- Edge cases
- Security considerations
- Performance aspects
