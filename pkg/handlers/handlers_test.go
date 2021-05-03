package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/storage"
)

var userid = ""

// helper function for mocking paths
func setup() *fiber.App {
	ctx := context.Background()
	app := fiber.New()
	NewHandler(ctx).ModifyPaths(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := setup()

	resp, err := app.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil {
		t.Fatalf("health check endpoint failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "ok" {
		t.Fatalf("health check endpoint failed: %v", err)
	}
}

func TestHandlerUserCreate(t *testing.T) {
	app := setup()
	firstPayload := map[string]interface{}{
		"should": "fail",
	}
	firstReqBody, _ := json.Marshal(firstPayload)

	req := httptest.NewRequest("POST", "/user/create", bytes.NewReader(firstReqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("user create endpoint failed failed: %v", err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body failed %v", err)
	}
	if resp.StatusCode != 422 {
		t.Fatalf("response should return status code: %d, got: %d", 422, resp.StatusCode)
	}

	secondPayload := map[string]interface{}{
		"display_name": "ycd_123",
		"country":      "us",
	}

	secondReqBody, _ := json.Marshal(secondPayload)

	req2 := httptest.NewRequest("POST", "/user/create", bytes.NewReader(secondReqBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("user create endpoint failed failed: %v", err)
	}
	defer resp2.Body.Close()

	var r2 struct {
		Data    storage.UserInfo `json:"data"`
		Err     string           `json:"error"`
		Success bool             `json:"success"`
	}

	err = json.NewDecoder(resp2.Body).Decode(&r2)
	if err != nil {
		t.Fatalf("unmarshalling response to struct failed: %v", err)
	}

	if r2.Success != true {
		t.Fatalf("request should be succeeded, but failed, %+v", r2)
	}

	userid = r2.Data.UserID
}

func TestHandlerScoreSubmit(t *testing.T) {
	app := setup()

	payload := map[string]interface{}{
		"score_worth": 55,
		"user_id":     userid,
	}

	reqBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/score/submit", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("user create endpoint failed failed: %v", err)
	}
	defer resp.Body.Close()

	var r Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatalf("unmarshalling response to struct failed: %v", err)
	}

	if r.Success != true {
		t.Fatalf("request should be succeeded, but failed, %+v", r)
	}
}

func TestHandleGetUser(t *testing.T) {
	app := setup()

	req := httptest.NewRequest("GET", "/user/profile/"+userid, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("user create endpoint failed failed: %v", err)
	}
	defer resp.Body.Close()

	var r Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatalf("unmarshalling response to struct failed: %v", err)
	}

	if r.Success != true {
		t.Fatalf("request should be succeeded, but failed, %+v", r)
	}
}

func TestHandlerLeaderboard(t *testing.T) {
	app := setup()

	req := httptest.NewRequest("GET", "/leaderboard", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("user create endpoint failed failed: %v", err)
	}
	defer resp.Body.Close()

	var r struct {
		Data    []storage.LeaderboardResult `json:"data"`
		Err     string                      `json:"error"`
		Success bool                        `json:"success"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatalf("unmarshalling response to struct failed: %v", err)
	}

	if r.Data[0].Country != "us" {
		t.Fatalf("expected country to be %s, got %s", "us", r.Data[0].Country)
	}
	if r.Data[0].Rank != 1 {
		t.Fatalf("expected rank to be %d, got %d", 1, r.Data[0].Rank)
	}
	if r.Data[0].DisplayName != "ycd_123" {
		t.Fatalf("expected rank to be %s, got %s", "ycd_123", r.Data[0].DisplayName)
	}

	if r.Success != true {
		t.Fatalf("request should be succeeded, but failed, %+v", r)
	}
}

func TestHandlerLeaderboardByCountry(t *testing.T) {
	app := setup()

	req := httptest.NewRequest("GET", "/leaderboard/us", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("user create endpoint failed failed: %v", err)
	}
	defer resp.Body.Close()

	var r struct {
		Data    []storage.LeaderboardResult `json:"data"`
		Err     string                      `json:"error"`
		Success bool                        `json:"success"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatalf("unmarshalling response to struct failed: %v", err)
	}

	if r.Data[0].Country != "us" {
		t.Fatalf("expected country to be %s, got %s", "us", r.Data[0].Country)
	}
	if r.Data[0].Rank != 1 {
		t.Fatalf("expected rank to be %d, got %d", 1, r.Data[0].Rank)
	}
	if r.Data[0].DisplayName != "ycd_123" {
		t.Fatalf("expected rank to be %s, got %s", "ycd_123", r.Data[0].DisplayName)
	}

	if r.Success != true {
		t.Fatalf("request should be succeeded, but failed, %+v", r)
	}
}
