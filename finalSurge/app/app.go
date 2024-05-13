package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	authUrl = "https://beta.finalsurge.com/api/login"
	activitiesUrl = "https://beta.finalsurge.com/api/WorkoutList"
)

type AuthResponse struct {
	ServerTime time.Time `json:"server_time"`
	Data       struct {
		Token         string      `json:"token"`
		UserKey       string      `json:"user_key"`
		FirstName     string      `json:"first_name"`
		LastName      string      `json:"last_name"`
		Email         string      `json:"email"`
		ProfilePicURL interface{} `json:"profile_pic_url"`
	} `json:"data"`
	Success          bool        `json:"success"`
	ErrorNumber      interface{} `json:"error_number"`
	ErrorDescription interface{} `json:"error_description"`
	CallID           interface{} `json:"call_id"`
}

type App struct {
	Email    string
	Password string
}

func NewApp(email, password string) *App {
	return &App{
		Email:    email,
		Password: password,
	}
}
// Authenticate the app created by the user
//
// Return the token and user key used for user auth
func (a *App) Authenticate(ctx context.Context) (*AuthResponse, error) {
	authPayload := map[string]string{"email": a.Email, "password": a.Password}
	jsonPayload, _ := json.Marshal(authPayload)

	req, err := http.NewRequestWithContext(ctx, "POST", authUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResp AuthResponse
	err1 := json.Unmarshal(responseData, &authResp)
	if err1 != nil {
		return nil, err1
	}
	return &authResp, nil
}
/*
func (a *App) GetActivities(ctx context.Context, scopeKey, userToken string, startDate, endDate time.Time) ([]byte, error) {
	activityPayload := map[string]string{
		"scope": "USER",
		"scopekey": scopeKey,
		"startdate": startDate.Format("2006-01-02"),
		"enddate": endDate.Format("2006-01-02"),
	}

	jsonPayload, _ := json.Marshal(activityPayload)
	req, err := http.NewRequestWithContext(ctx, "GET", activitiesUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
*/


func (a *App) GetActivities(ctx context.Context, userToken, scopeKey string, startDate, endDate time.Time) ([]byte, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", userToken),
	}

	params := map[string]string{
		"scope":     "USER",
		"scopekey":  scopeKey,
		"startdate": startDate.Format("2006-01-02"),
		"enddate":   endDate.Format("2006-01-02"),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", activitiesUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set query parameters
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}