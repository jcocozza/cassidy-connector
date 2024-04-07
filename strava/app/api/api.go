package api

import (
	"context"
	"time"

	"github.com/antihax/optional"
	"github.com/jcocozza/cassidy-connector/strava/internal/swagger"
)

type StravaAPI struct {
	stravaClient *swagger.APIClient
}
func NewStravaAPI(stravaClient *swagger.APIClient) *StravaAPI {
	return &StravaAPI{
		stravaClient: stravaClient,
	}
}
// Get the athlete that is logged-in/authenticated
func (api *StravaAPI) GetAthlete(ctx context.Context) (*swagger.DetailedAthlete, error) {
    athlete, _, err := api.stravaClient.AthletesApi.GetLoggedInAthlete(ctx)
    if err != nil {
        return nil, err
    }
    return &athlete, nil
}
// Get activities. Will enumerate through all available pages of data.
//
// before and after are times to filter activies by. Both are optional (pass in nil to ignore them)
// 	- before will filter for activities before a passed time.Time
//	- after will filter for activities after the passed time.Time
// before and after are converted to epoch timestamp integers.
//
// perPage is the number of activities per page. (default 30) (max 200)
//
// If you plan on retreiving lots of data, you should set per page to be high. This will drastically reduce the number of API calls made.
// (There is an API call made for each page)
func (api *StravaAPI) GetActivities(ctx context.Context, perPage int, before, after *time.Time) ([][]swagger.SummaryActivity, error) {
	var summaryActivitylol [][]swagger.SummaryActivity
	opts := &swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts{}
	if before != nil {
		beforeOpt := optional.NewInt32(int32(before.Unix()))
		opts.Before = beforeOpt
	}
	if after != nil {
		afterOpt := optional.NewInt32(int32(after.Unix()))
		opts.After = afterOpt
	}
	perPageOpt := optional.NewInt32(int32(perPage))
	opts.PerPage = perPageOpt

	existsMore := true
	var page int32 = 1 // page enumeration starts at 1
	for existsMore { // enumerate until there are no more activities
		opts.Page = optional.NewInt32(page)
		summary, _, err := api.stravaClient.ActivitiesApi.GetLoggedInAthleteActivities(ctx, opts)
		if err != nil {
			return nil, err
		}
		//return summary, nil
		summaryActivitylol = append(summaryActivitylol, summary)
		page += 1

		if len(summary) == 0 {
			existsMore = false
		}
	}
	return summaryActivitylol, nil
}