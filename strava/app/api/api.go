package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"github.com/jcocozza/cassidy-connector/strava/swagger"
	"golang.org/x/oauth2"
)

// if the strava api returns 404 not found, will throw this error
var NotFoundError = errors.New("Object not found")

// StreamType represents the different types of steams that exist.
// These are exported from the package.
type StreamType string

const (
	Time           StreamType = "time"            // time stream
	Distance       StreamType = "distance"        // distance stream
	Latlng         StreamType = "latlng"          // latlng stream
	Altitude       StreamType = "altitude"        // altitude stream
	VelocitySmooth StreamType = "velocity_smooth" // velocity stream
	Heartrate      StreamType = "heartrate"       // heartrate stream
	Cadence        StreamType = "cadence"         // cadence stream
	Watts          StreamType = "watts"           // watts stream
	Temp           StreamType = "temp"            // temp stream
	Moving         StreamType = "moving"          // moving stream
	GradeSmooth    StreamType = "grade_smooth"    // grade stream
)

// This contains the user's short-lived access token which is used to access data.
// When it expires, use the user's refresh token to get a new access token.
//
// This struct is obtained in 1 of 2 ways:
//   - by possessing an existing refresh token and getting a new access token (handled automatically by the oauth2 package).
//   - via user authorization, whereby an auth code is issued and is used to get the access token.
type userSession struct {
	tkn *oauth2.Token
}

// this is needed to properly work with the swagger implementation
//
// when doing the OAuth, swagger will call this method to get the token
func (us *userSession) Token() (*oauth2.Token, error) {
	return us.tkn, nil
}

// return context for proper authorization when sending to the api
func (us *userSession) AuthorizationContext(parent context.Context) context.Context {
	return context.WithValue(parent, swagger.ContextOAuth2, us)
}

// The StravaAPI struct is the primary means of interacting with the strava api.
//
// This is the layer of abstraction so that users don't have to directly deal with api calls.
//
// It will handle methods related to data in the strava api. (auth will be handled by the broader app struct)
//
// Whenever possible, this will throw the NotFoundError when the underlying strava api returns a 404
type StravaAPI struct {
	stravaClient *swagger.APIClient
}

func NewStravaAPI(stravaClient *swagger.APIClient) *StravaAPI {
	return &StravaAPI{
		stravaClient: stravaClient,
	}
}

// Get the athlete that is logged-in/authenticated
func (api *StravaAPI) GetAthlete(ctx context.Context, token *oauth2.Token) (*swagger.DetailedAthlete, error) {
	us := &userSession{tkn: token}
	ctx = us.AuthorizationContext(ctx)
	athlete, resp, err := api.stravaClient.AthletesApi.GetLoggedInAthlete(ctx)
	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}
	if err != nil {
		return nil, err
	}
	return &athlete, nil
}

// Get activities. Will cycle through all available pages of data.
//
// before and after are times to filter activies by. Both are optional (pass in nil to ignore them)
//   - before will filter for activities before a passed time.Time
//   - after will filter for activities after the passed time.Time
//
// before and after are converted to epoch timestamp integers.
//
// perPage is the number of activities per page. (default 30) (max 200)
//
// If you plan on retreiving lots of data, you should set per page to be high. This will drastically reduce the number of API calls made.
// (There is an API call made for each page)
func (api *StravaAPI) GetActivities(ctx context.Context, token *oauth2.Token, perPage int, before, after *time.Time) ([][]swagger.SummaryActivity, error) {
	us := &userSession{tkn: token}
	ctx = us.AuthorizationContext(ctx)
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
	for existsMore {   // enumerate until there are no more activities
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
func (api *StravaAPI) GetActivity(ctx context.Context, token *oauth2.Token, activityID int, includeAllEfforts bool) (*swagger.DetailedActivity, error) {
	us := &userSession{tkn: token}
	ctx = us.AuthorizationContext(ctx)
	opts := &swagger.ActivitiesApiGetActivityByIdOpts{IncludeAllEfforts: optional.NewBool(includeAllEfforts)}
	activity, resp, err := api.stravaClient.ActivitiesApi.GetActivityById(ctx, int64(activityID), opts)
	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}
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
//   - time, distance, latlng, altitude, velocity_smooth, heartrate, cadence, watts, temp, moving, grade_smooth
//
// Each of these are exported by the package as constant symbols for proper access as StreamType types.
//
// This will return a struct containing the desired streams for the activity
func (api *StravaAPI) GetActivityStreams(ctx context.Context, token *oauth2.Token, activityID int, keys []StreamType) (*swagger.StreamSet, error) {
	us := &userSession{tkn: token}
	ctx = us.AuthorizationContext(ctx)
	keyByType := true
	err := validateKeys(keys)
	if err != nil {
		return nil, err
	}
	keyList := convertKeys(keys)
	streamSet, resp, err := api.stravaClient.StreamsApi.GetActivityStreams(ctx, int64(activityID), keyList, keyByType)
	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}
	if err != nil {
		return nil, err
	}
	return &streamSet, nil
}
