package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	main "github.com/obedtandadjaja/project_k_backend/services/auth"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
)

var app main.App

func TestMain(m *testing.M) {
	app = main.App{}
	app.Initialize(
		os.Getenv("ENV"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	code := m.Run()

	os.Exit(code)
}

func clearCredentialsTable() {
	app.DB.Exec("delete from sessions")
	app.DB.Exec("delete from credentials")
}

func executeRequest(req *http.Request) (*httptest.ResponseRecorder, map[string]interface{}) {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	var responseBody map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &responseBody)

	return rr, responseBody
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateCredentialInvalidRequest(t *testing.T) {
	payload := []byte(`{"credential_id":0,"password":0}`)

	req, _ := http.NewRequest("POST", "/credentials", bytes.NewBuffer(payload))
	rr, _ := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, rr.Code)
}

func TestCreateCredential(t *testing.T) {
	clearCredentialsTable()

	rr, responseBody := createCredential("email@email.com", "5591111111", "password")
	checkResponseCode(t, http.StatusCreated, rr.Code)

	if uuid, ok := responseBody["credential_uuid"]; ok {
		credentials, _ := credential.All(app.DB)

		if len(credentials) == 1 {
			c := credentials[0]

			if c.Uuid != uuid {
				t.Errorf("Created credential id is wrong")
			}
		} else {
			t.Errorf("Expected one credential to be created, found %d", len(credentials))
		}
	} else {
		t.Errorf("Missing uuid in response")
	}
}

func createCredential(email, phone, password string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"email":"%s","phone":"%s","password":"%s"}`, email, phone, password)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/credentials", bytes.NewBuffer(payload))
	rr, responseBody := executeRequest(req)

	return rr, responseBody
}

// Consider removing delete credential, since there is no use case for it. If we do want accounts
// to be deactivated, it should be a soft delete instead
func TestDeleteCredential(t *testing.T) {
	clearCredentialsTable()

	rr, createResponseBody := createCredential("email@email.com", "5591111111", "password")

	rr, deleteResponseBody := deleteCredential(createResponseBody["credential_uuid"].(string))

	if createResponseBody["credential_uuid"] != deleteResponseBody["credential_uuid"] {
		t.Errorf("Uuid for creation and deletion do not match")
	}

	checkResponseCode(t, http.StatusNoContent, rr.Code)
}

func deleteCredential(credentialId string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"credential_uuid":"%s"}`, credentialId)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("DELETE", "/credentials", bytes.NewBuffer(payload))
	rr, responseBody := executeRequest(req)

	return rr, responseBody
}

// Authentication flow
// ------------------------------------
// 1. create credential
// 2. login - get session
// 2a. verify that access token works
// 2b. verify that refresh token works
// 3. token - request new access token
// 3a. verify that access token works

func TestLogin(t *testing.T) {
	clearCredentialsTable()

	rr, createResponseBody := createCredential("email@email.com", "", "password")

	// try logging in with cred uuid
	rr, loginResponseBody := loginWithUuid(createResponseBody["credential_uuid"].(string), "password")
	checkResponseCode(t, http.StatusOK, rr.Code)

	// try logging in with email
	rr, loginResponseBody = loginWithEmail("email@email.com", "password")
	checkResponseCode(t, http.StatusOK, rr.Code)

	jwt := loginResponseBody["jwt"].(string)
	sessionJwt := loginResponseBody["session"].(string)

	rr, verifyResponseBody := verifyToken(jwt)
	checkResponseCode(t, http.StatusOK, rr.Code)

	if !verifyResponseBody["verified"].(bool) {
		t.Errorf("Invalid jwt token")
	}

	rr, tokenResponseBody := token(sessionJwt)
	checkResponseCode(t, http.StatusOK, rr.Code)

	if _, ok := tokenResponseBody["jwt"]; !ok {
		t.Errorf("Invalid session jwt")
	}
}

func TestTokenAfterCreate(t *testing.T) {
	clearCredentialsTable()

	rr, createResponseBody := createCredential("email@email.com", "", "password")
	checkResponseCode(t, http.StatusCreated, rr.Code)

	jwt := createResponseBody["jwt"].(string)
	session := createResponseBody["session"].(string)

	rr, verifyResponseBody := verifyToken(jwt)
	checkResponseCode(t, http.StatusOK, rr.Code)

	if !verifyResponseBody["verified"].(bool) {
		t.Errorf("Invalid jwt token")
	}

	rr, tokenResponseBody := token(session)
	checkResponseCode(t, http.StatusOK, rr.Code)

	if _, ok := tokenResponseBody["jwt"]; !ok {
		t.Errorf("Invalid session jwt")
	}
}

func TestVerifySessionToken(t *testing.T) {
	clearCredentialsTable()

	rr, createResponseBody := createCredential("email@email.com", "", "password")
	checkResponseCode(t, http.StatusCreated, rr.Code)

	session := createResponseBody["session"].(string)

	rr, _ = verifySessionToken(session)
	checkResponseCode(t, http.StatusOK, rr.Code)
}

func loginWithUuid(credentialUuid, password string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"credential_uuid":"%s","password":"%s"}`, credentialUuid, password)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	rr, responseBody := executeRequest(req)

	return rr, responseBody
}

func loginWithEmail(email, password string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	rr, responseBody := executeRequest(req)

	return rr, responseBody
}

func token(sessionJwt string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"session":"%s"}`, sessionJwt)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/token", bytes.NewBuffer(payload))

	rr, responseBody := executeRequest(req)

	return rr, responseBody
}

func verifyToken(accessTokenJwt string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"jwt":"%s"}`, accessTokenJwt)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(payload))

	rr, responseBody := executeRequest(req)

	return rr, responseBody
}

func verifySessionToken(sessionTokenJwt string) (*httptest.ResponseRecorder, map[string]interface{}) {
	jsonString := fmt.Sprintf(`{"session":"%s"}`, sessionTokenJwt)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/verify_session_token", bytes.NewBuffer(payload))

	rr, responseBody := executeRequest(req)

	return rr, responseBody
}
