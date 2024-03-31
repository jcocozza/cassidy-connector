package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)
const (
	athleteUrl string = "https://www.strava.com/api/v3/athlete"
	activitiesUrl string = "https://www.strava.com/api/v3/athlete/activities"
)
type StravaAPI interface {
	// Get activities by page
	//
	// numPages determines the number of pages to get.
	// perPage is the number of activities per page. (max 200)
	GetActivityPages(accessToken string, numPages int, perPage int) ([][]byte, error)
	GetAthlete(accessToken string)
}

type StravaAPICaller struct {}
// Get activities by page
//
// numPages determines the number of pages to get.
// perPage is the number of activities per page. (max 200)
func (sac *StravaAPICaller) GetActivityPages(accessToken string, numPages int, perPage int) ([][]byte, error) {
	pages := [][]byte{}
	// page enumerate starts at 1
	for i := 1; i <= numPages; i++ {
		query := url.Values{}
		query.Set("per_page", fmt.Sprint(perPage))
		query.Set("page", fmt.Sprint(i))

		urlWithParams := fmt.Sprintf("%s?%s", activitiesUrl, query.Encode())

		// Create a new GET request
		req, err := http.NewRequest("GET", urlWithParams, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return nil, err
		}
		// Set Authorization header with access token
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		// Perform the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error performing request:", err)
			return nil, err
		}
		defer resp.Body.Close()
		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil, err
		}
		// the api doesn't quite return an empty body I guess?
		if len(body) == 2 {
			break
		}
		pages = append(pages, body)
	}
	return pages, nil
}
func (sac *StravaAPICaller) GetAthlete(accessToken string) {
	// Create a new GET request
	req, err := http.NewRequest("GET", athleteUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set Authorization header with access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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