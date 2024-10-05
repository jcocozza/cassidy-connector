package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"github.com/jcocozza/cassidy-connector/strava/swagger"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

// if the strava api returns 404 not found, will throw this error
var NotFoundError = errors.New("Object not found")

// if the rate limiter throws an error
var RateLimitError = errors.New("Rate Limit Error. (this likely means context expired while waiting for rate limits to be reset)")

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

const (
	// The strava API limits to 300 READ requests per 15 minutes
	ReadRateLimit15Min = 300
	// The strava API limits to 3000 READ requests per day
	ReadRateLimiteDaily = 3000
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
// It will handle methods related to data in the strava api.
//
// # Whenever possible, this will return the NotFoundError when the underlying strava api returns a 404
//
// All the methods here are also rate limited per the strava guidelines
// Make sure that every context has a timeout, otherwise the program will block until the rate limits refreshes.
type StravaAPI struct {
	stravaClient *swagger.APIClient
	logger       *slog.Logger
	oauth        *oauth2.Config
	limiter15min *rate.Limiter
	limiterDaily *rate.Limiter
}

func NewStravaAPI(stravaClient *swagger.APIClient, cfg *oauth2.Config, logger *slog.Logger) *StravaAPI {
	return &StravaAPI{
		stravaClient: stravaClient,
		logger:       logger,
		oauth:        cfg,
		limiter15min: rate.NewLimiter(rate.Every(15*time.Minute), ReadRateLimit15Min),
		limiterDaily: rate.NewLimiter(rate.Every(24*time.Hour), ReadRateLimiteDaily),
	}
}

// check to see if the limits have been surpassed
//
// ** should be called before every api call **
func (api *StravaAPI) checkRateLimits(ctx context.Context) error {
	err := api.limiterDaily.Wait(ctx)
	if err != nil {
		api.logger.ErrorContext(ctx, "failed daily rate limits", slog.String("error", err.Error()))
		return RateLimitError
	}
	err = api.limiter15min.Wait(ctx)
	if err != nil {
		api.logger.ErrorContext(ctx, "failed 15 minute rate limits", slog.String("error", err.Error()))
		return RateLimitError
	}
	return nil
}

// return the remaining requests for the 15 mintue request window and the daily window
// (in that order)
func (api *StravaAPI) RemainingRequests() (int, int) {
	rr15 := int(api.limiter15min.Tokens())
	rrdaily := int(api.limiterDaily.Tokens())
	return rr15, rrdaily
}

// auto refresh the token via TokenSource
func (api *StravaAPI) refreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	src := api.oauth.TokenSource(ctx, token)
	newToken, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return newToken, nil
}

// set auth context with a refreshed token
func (api *StravaAPI) setContext(ctx context.Context, token *oauth2.Token) (context.Context, error) {
	refreshedTkn, err := api.refreshToken(ctx, token)
	if err != nil {
		api.logger.ErrorContext(ctx, "token refresh failed")
		return nil, err
	}
	us := &userSession{tkn: refreshedTkn}
	newCtx := us.AuthorizationContext(ctx)
	return newCtx, nil
}

// Get the athlete that is logged-in/authenticated
func (api *StravaAPI) GetAthlete(ctx context.Context, token *oauth2.Token) (*swagger.DetailedAthlete, error) {
	err := api.checkRateLimits(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = api.setContext(ctx, token)
	if err != nil {
		return nil, err
	}
	api.logger.DebugContext(ctx, "getting athlete")
	athlete, resp, err := api.stravaClient.AthletesApi.GetLoggedInAthlete(ctx)
	if resp.StatusCode == http.StatusNotFound {
		api.logger.DebugContext(ctx, "athlete not found")
		return nil, NotFoundError
	}
	if err != nil {
		api.logger.ErrorContext(ctx, "error getting athlete", slog.String("error", err.Error()))
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
	ctx, err := api.setContext(ctx, token)
	if err != nil {
		return nil, err
	}
	api.logger.DebugContext(ctx, "getting activities",
		slog.Int("per page", perPage),
		slog.Any("before", before),
		slog.Any("after", after),
	)
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
		err := api.checkRateLimits(ctx)
		if err != nil {
			return nil, err
		}
		summary, _, err := api.stravaClient.ActivitiesApi.GetLoggedInAthleteActivities(ctx, opts)
		if err != nil {
			api.logger.ErrorContext(ctx, "getting activities failed", slog.Any("page", opts.Page), slog.String("error", err.Error()))
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
	err := api.checkRateLimits(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = api.setContext(ctx, token)
	if err != nil {
		return nil, err
	}
	api.logger.DebugContext(ctx, "getting activity", slog.Int("activity id", activityID), slog.Bool("include all efforts", includeAllEfforts))
	opts := &swagger.ActivitiesApiGetActivityByIdOpts{IncludeAllEfforts: optional.NewBool(includeAllEfforts)}
	activity, resp, err := api.stravaClient.ActivitiesApi.GetActivityById(ctx, int64(activityID), opts)
	if resp.StatusCode == http.StatusNotFound {
		api.logger.DebugContext(ctx, "activity not found")
		return nil, NotFoundError
	}
	if err != nil {
		api.logger.ErrorContext(ctx, "error getting activity", slog.String("error", err.Error()))
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
	ctx, err := api.setContext(ctx, token)
	if err != nil {
		return nil, err
	}
	api.logger.DebugContext(ctx, "getting activity streams", slog.Int("activity id", activityID), slog.Any("keys", keys))
	keyByType := true
	err = validateKeys(keys)
	if err != nil {
		api.logger.ErrorContext(ctx, "invalid keys", slog.String("error", err.Error()))
		return nil, err
	}
	keyList := convertKeys(keys)
	err = api.checkRateLimits(ctx)
	if err != nil {
		return nil, err
	}
	streamSet, resp, err := api.stravaClient.StreamsApi.GetActivityStreams(ctx, int64(activityID), keyList, keyByType)
	if resp.StatusCode == http.StatusNotFound {
		api.logger.DebugContext(ctx, "streams not found")
		return nil, NotFoundError
	}
	if err != nil {
		api.logger.ErrorContext(ctx, "error getting streams", slog.String("error", err.Error()))
		return nil, err
	}
	return &streamSet, nil
}
