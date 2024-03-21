package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jcocozza/cassidy-connector/strava/auth"
)

const (
	AthleteUrl string = "https://www.strava.com/api/v3/athlete"
)

func GetAthlete(refreshToken string) {
	tokenResp, err := auth.GetAccessTokenFromRefresh(refreshToken)
	if err != nil {
		return
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", AthleteUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set Authorization header with access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenResp.AccessToken))

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print response body
	fmt.Println(string(body))
}

func main() {
	//GetAthlete()
	//GetActivities()
}