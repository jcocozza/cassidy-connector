package api

import (
	"context"
	"fmt"
	"time"

	"github.com/antihax/optional"
	"github.com/jcocozza/cassidy-connector/strava/internal/swagger"
)

// StreamType represents the different types of steams that exist.
// These are exported from the package.
type StreamType string
const (
	Time StreamType = "time" // time stream
	Distance StreamType = "distance" // distance stream
	Latlng StreamType = "latlng" // latlng stream
	Altitude StreamType = "altitude" // altitude stream
	VelocitySmooth StreamType = "velocity_smooth" // velocity stream
	Heartrate StreamType = "heartrate" // heartrate stream
	Cadence StreamType = "cadence" // cadence stream
	Watts StreamType = "watts" // watts stream
	Temp StreamType = "temp" // temp stream
	Moving StreamType = "moving" // moving stream
	GradeSmooth StreamType = "grade_smooth" // grade stream
)
// The StravaAPI struct is the primary means of interacting with the strava api.
//
// This is the layer of abstraction so that users don't have to directly deal with api calls.
//
// It will handle methods related to data in the strava api. (auth will be handled by the broader app struct)
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
// Get activities. Will cycle through all available pages of data.
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
// Get a single activity by activity ID
//
// `activityID` is the id of the activity
//
// `includeAllEfforts` includes all segment efforts if true
func (api *StravaAPI) GetActivity(ctx context.Context, activityID int, includeAllEfforts bool) (*swagger.DetailedActivity, error) {
	opts := &swagger.ActivitiesApiGetActivityByIdOpts{IncludeAllEfforts: optional.NewBool(includeAllEfforts)}
	activity, _, err := api.stravaClient.ActivitiesApi.GetActivityById(ctx, int64(activityID), opts)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}
// convert a list of StreamType into a list of string
//
// this is just a simple way to ensure that users aren't passing weird stream types into the `GetActivityStreams` function
func convertKeys(keys []StreamType) []string {
	l := []string{}
	for _, key := range keys {
		l = append(l, string(key))
	}
	return l
}
func validateKeys(keys []StreamType) error {
	keyList := []StreamType{Time, Distance, Latlng, Altitude, VelocitySmooth, Heartrate, Cadence, Watts, Temp, Moving, GradeSmooth}

	for _, key := range keys {
		isInvalid := true
		for _, actKey := range keyList {
			if actKey == key {
				isInvalid = false
				break
			}
		}
		if isInvalid {
			return fmt.Errorf("%s is not a correct stream type", key)
		}
	}
	return nil
}
// Get the streams for a given activity.
//
// `activityID` is the id of the activity
//
// `keys` is a list of the kinds of streams you want to get for that activity.
// Currently, the following keys(stream types) are supported by strava:
//		- time, distance, latlng, altitude, velocity_smooth, heartrate, cadence, watts, temp, moving, grade_smooth
// Each of these are exported by the package as constant symbols for proper access as StreamType types.
//
// This will return a struct containing the desired streams for the activity
func (api *StravaAPI) GetActivityStreams(ctx context.Context, activityID int, keys []StreamType) (*swagger.StreamSet, error) {
	keyByType := true
	err := validateKeys(keys)
	if err != nil {
		return nil, err
	}
	keyList := convertKeys(keys)
	streamSet, _, err := api.stravaClient.StreamsApi.GetActivityStreams(ctx, int64(activityID), keyList, keyByType)
	if err != nil {
		return nil, err
	}
	return &streamSet, nil
}