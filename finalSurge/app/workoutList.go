package app

import "time"

// Used https://transform.tools/json-to-go to convert a json response to struct
type WorkoutListResponse struct {
	ServerTime time.Time `json:"server_time"`
	Data       []struct {
		HasDownloadFile        bool        `json:"has_download_file"`
		DownloadFileExtension  string      `json:"download_file_extension"`
		AppleSyncTime          interface{} `json:"apple_sync_time"`
		AppleSyncUUID          interface{} `json:"apple_sync_uuid"`
		WahooSyncTime          interface{} `json:"wahoo_sync_time"`
		WahooSyncAvailable     bool        `json:"wahoo_sync_available"`
		GarminSyncTime         interface{} `json:"garmin_sync_time"`
		GarminSyncAvailable    bool        `json:"garmin_sync_available"`
		ExternalDataSource     string      `json:"external_data_source"`
		IsTeamWorkout          bool        `json:"is_team_workout"`
		CanSplit               bool        `json:"can_split"`
		ValidMergeTarget       bool        `json:"valid_merge_target"`
		ValidMergeSource       bool        `json:"valid_merge_source"`
		HasPainPointRecords    bool        `json:"has_pain_point_records"`
		HasStructuredWorkout   bool        `json:"has_structured_workout"`
		PainPointRecords       interface{} `json:"pain_point_records"`
		Gender                 string      `json:"gender"`
		UserKey                string      `json:"user_key"`
		UserName               string      `json:"user_name"`
		UserProfilePicURL      interface{} `json:"user_profile_pic_url"`
		Key                    string      `json:"key"`
		WcalKey                interface{} `json:"wcal_key"`
		WcalLabel              interface{} `json:"wcal_label"`
		CanDelete              bool        `json:"can_delete"`
		CanHide                bool        `json:"can_hide"`
		CanMove                bool        `json:"can_move"`
		CanEdit                bool        `json:"can_edit"`
		CoachAssigned          bool        `json:"coach_assigned"`
		CoachUserKey           interface{} `json:"coach_user_key"`
		CoachName              interface{} `json:"coach_name"`
		CoachProfilePicURL     interface{} `json:"coach_profile_pic_url"`
		HasActualData          bool        `json:"has_actual_data"`
		HasIntervals           bool        `json:"has_intervals"`
		HasMap                 bool        `json:"has_map"`
		HasStats               bool        `json:"has_stats"`
		HasAttachments         bool        `json:"has_attachments"`
		HasRoutes              bool        `json:"has_routes"`
		Attachments            interface{} `json:"attachments"`
		WorkoutDate            string      `json:"workout_date"`
		WorkoutTime            string      `json:"workout_time"`
		Order                  int         `json:"order"`
		PlanDay                interface{} `json:"plan_day"`
		Name                   string      `json:"name"`
		Description            string      `json:"description"`
		LocationName           interface{} `json:"location_name"`
		LocationStreet         interface{} `json:"location_street"`
		LocationCity           interface{} `json:"location_city"`
		LocationState          interface{} `json:"location_state"`
		LocationZip            interface{} `json:"location_zip"`
		LocationCountry        interface{} `json:"location_country"`
		IsRace                 bool        `json:"is_race"`
		RacePlaceOverall       interface{} `json:"race_place_overall"`
		RaceAgeGroup           interface{} `json:"race_age_group"`
		Felt                   interface{} `json:"felt"`
		Effort                 interface{} `json:"effort"`
		PostWorkoutNotes       interface{} `json:"post_workout_notes"`
		WeatherTemperature     interface{} `json:"weather_temperature"`
		WeatherIsCelsius       interface{} `json:"weather_is_celsius"`
		WeatherHumidity        interface{} `json:"weather_humidity"`
		WeatherSunny           interface{} `json:"weather_sunny"`
		WeatherPartlySunny     interface{} `json:"weather_partly_sunny"`
		WeatherCloudy          interface{} `json:"weather_cloudy"`
		WeatherLightrain       interface{} `json:"weather_lightrain"`
		WeatherRain            interface{} `json:"weather_rain"`
		WeatherSnow            interface{} `json:"weather_snow"`
		WeatherWindy           interface{} `json:"weather_windy"`
		CommentCount           int         `json:"CommentCount"`
		CommentCountNew        interface{} `json:"CommentCountNew"`
		WorkoutCompletion      int         `json:"workout_completion"`
		WorkoutStatusText      string      `json:"workout_status_text"`
		WorkoutStatusColor     string      `json:"workout_status_color"`
		WorkoutStatusIndicator int         `json:"workout_status_indicator"`
		MapURL                 string      `json:"MapURL"`
		Activities             []struct {
			ActivityTypeKey       string      `json:"activity_type_key"`
			ActivityTypeName      string      `json:"activity_type_name"`
			ActivityTypeIcon      int         `json:"activity_type_icon"`
			ActivitySubTypeKey    interface{} `json:"activity_sub_type_key"`
			ActivitySubTypeName   interface{} `json:"activity_sub_type_name"`
			ActivityTypeColor     string      `json:"activity_type_color"`
			ActivityTypeForecolor string      `json:"activity_type_forecolor"`
			Equipment             struct {
				EquipmentKey    string      `json:"equipment_key"`
				EquipmentTypeID int         `json:"equipment_type_id"`
				EquipmentName   string      `json:"equipment_name"`
				EquipmentNotes  interface{} `json:"equipment_notes"`
				EquipmentBrand  struct {
					BrandKey  string `json:"brand_key"`
					BrandName string `json:"brand_name"`
				} `json:"equipment_brand"`
				EquipmentModel                   string      `json:"equipment_model"`
				EquipmentCost                    interface{} `json:"equipment_cost"`
				EquipmentPurchasedate            interface{} `json:"equipment_purchasedate"`
				EquipmentRetiredate              interface{} `json:"equipment_retiredate"`
				EquipmentDistance                float64     `json:"equipment_distance"`
				EquipmentStartDistance           float64     `json:"equipment_start_distance"`
				EquipmentStartDistanceUnit       string      `json:"equipment_start_distance_unit"`
				EquipmentAlertDistanceNormalized float64     `json:"equipment_alert_distance_normalized"`
				EquipmentAlertDistance           interface{} `json:"equipment_alert_distance"`
				EquipmentAlertDistanceUnit       string      `json:"equipment_alert_distance_unit"`
			} `json:"equipment"`
			Route                    interface{} `json:"route"`
			Number                   int         `json:"number"`
			PlannedDuration          float64     `json:"planned_duration"`
			PlannedAmount            float64     `json:"planned_amount"`
			PlannedAmountType        string      `json:"planned_amount_type"`
			PlannedAmountNormalized  float64     `json:"planned_amount_normalized"`
			PlannedPaceLow           interface{} `json:"planned_pace_low"`
			PlannedPaceLowType       interface{} `json:"planned_pace_low_type"`
			PlannedPaceHigh          interface{} `json:"planned_pace_high"`
			PlannedPaceHighType      interface{} `json:"planned_pace_high_type"`
			PlannedPaceDisplay       interface{} `json:"planned_pace_display"`
			PlannedPaceDisplayType   interface{} `json:"planned_pace_display_type"`
			Quantity                 int         `json:"quantity"`
			Duration                 float64     `json:"duration"`
			TimeElapsed              float64     `json:"time_elapsed"`
			TimeTimer                float64     `json:"time_timer"`
			TimeMoving               float64     `json:"time_moving"`
			Amount                   float64     `json:"amount"`
			AmountType               string      `json:"amount_type"`
			AmountNormalized         float64     `json:"amount_normalized"`
			Pace                     float64     `json:"pace"`
			PaceType                 string      `json:"pace_type"`
			PaceDisplay              string      `json:"pace_display"`
			PaceDisplayType          string      `json:"pace_display_type"`
			SpeedAvg                 float64     `json:"speed_avg"`
			SpeedMax                 float64     `json:"speed_max"`
			SpeedType                string      `json:"speed_type"`
			TempAvg                  interface{} `json:"temp_avg"`
			TempMax                  interface{} `json:"temp_max"`
			PowerAvg                 int         `json:"power_avg"`
			PowerMax                 int         `json:"power_max"`
			CadenceAvg               int         `json:"cadence_avg"`
			CadenceMax               int         `json:"cadence_max"`
			HrAvg                    int         `json:"hr_avg"`
			HrMax                    int         `json:"hr_max"`
			RpmAvg                   interface{} `json:"rpm_avg"`
			RpmMax                   interface{} `json:"rpm_max"`
			ElevationGainDisplayType string      `json:"elevation_gain_display_type"`
			ElevationGainDisplay     string      `json:"elevation_gain_display"`
			ElevationLossDisplayType string      `json:"elevation_loss_display_type"`
			ElevationLossDisplay     string      `json:"elevation_loss_display"`
			ElevationGain            float64     `json:"elevation_gain"`
			ElevationGainType        string      `json:"elevation_gain_type"`
			ElevationLoss            float64     `json:"elevation_loss"`
			ElevationLossType        string      `json:"elevation_loss_type"`
			Calories                 int         `json:"calories"`
			Variability              float64     `json:"variability"`
			Intensity                interface{} `json:"intensity"`
			WeightedPower            int         `json:"weighted_power"`
			MeanmaxPower30           int         `json:"meanmax_power_30"`
			VerticalOscillationAvg   interface{} `json:"vertical_oscillation_avg"`
			VerticalOscillationMax   interface{} `json:"vertical_oscillation_max"`
			GroundContactTimeAvg     interface{} `json:"ground_contact_time_avg"`
			GroundContactTimeMax     interface{} `json:"ground_contact_time_max"`
			GroundContactBalanceAvg  interface{} `json:"ground_contact_balance_avg"`
			GroundContactBalanceMax  interface{} `json:"ground_contact_balance_max"`
			StrideLengthAvg          float64     `json:"stride_length_avg"`
			VerticalRatioAvg         interface{} `json:"vertical_ratio_avg"`
			FormPower                interface{} `json:"form_power"`
			LegSpring                interface{} `json:"leg_spring"`
			RightPowerAvg            interface{} `json:"right_power_avg"`
			RightPowerPctAvg         interface{} `json:"right_power_pct_avg"`
			LeftPowerAvg             interface{} `json:"left_power_avg"`
			LeftPowerPctAvg          interface{} `json:"left_power_pct_avg"`
			RestActivity             interface{} `json:"RestActivity"`
			Laps                     []struct {
				Number                   int         `json:"number"`
				Quantity                 interface{} `json:"quantity"`
				Duration                 float64     `json:"duration"`
				Amount                   float64     `json:"amount"`
				AmountType               string      `json:"amount_type"`
				AmountNormalized         float64     `json:"amount_normalized"`
				Pace                     interface{} `json:"pace"`
				PaceType                 interface{} `json:"pace_type"`
				PaceDisplay              string      `json:"pace_display"`
				PaceDisplayType          string      `json:"pace_display_type"`
				SpeedAvg                 float64     `json:"speed_avg"`
				SpeedMax                 float64     `json:"speed_max"`
				SpeedType                string      `json:"speed_type"`
				TempAvg                  int         `json:"temp_avg"`
				TempMax                  interface{} `json:"temp_max"`
				PowerAvg                 int         `json:"power_avg"`
				PowerMax                 int         `json:"power_max"`
				CadenceAvg               int         `json:"cadence_avg"`
				CadenceMax               int         `json:"cadence_max"`
				HrAvg                    int         `json:"hr_avg"`
				HrMax                    int         `json:"hr_max"`
				RpmAvg                   interface{} `json:"rpm_avg"`
				RpmMax                   interface{} `json:"rpm_max"`
				ElevationGainDisplayType string      `json:"elevation_gain_display_type"`
				ElevationGainDisplay     string      `json:"elevation_gain_display"`
				ElevationLossDisplayType string      `json:"elevation_loss_display_type"`
				ElevationLossDisplay     string      `json:"elevation_loss_display"`
				ElevationGain            float64     `json:"elevation_gain"`
				ElevationGainType        string      `json:"elevation_gain_type"`
				ElevationLoss            float64     `json:"elevation_loss"`
				ElevationLossType        string      `json:"elevation_loss_type"`
				Calories                 int         `json:"calories"`
				VerticalOscillationAvg   int         `json:"vertical_oscillation_avg"`
				VerticalOscillationMax   interface{} `json:"vertical_oscillation_max"`
				GroundContactTimeAvg     float64     `json:"ground_contact_time_avg"`
				GroundContactTimeMax     interface{} `json:"ground_contact_time_max"`
				GroundContactBalanceAvg  interface{} `json:"ground_contact_balance_avg"`
				GroundContactBalanceMax  interface{} `json:"ground_contact_balance_max"`
				StrideLengthAvg          float64     `json:"stride_length_avg"`
				VerticalRatioAvg         float64     `json:"vertical_ratio_avg"`
				FormPower                interface{} `json:"form_power"`
				LegSpring                interface{} `json:"leg_spring"`
				RightPowerAvg            interface{} `json:"right_power_avg"`
				RightPowerPctAvg         interface{} `json:"right_power_pct_avg"`
				LeftPowerAvg             interface{} `json:"left_power_avg"`
				LeftPowerPctAvg          interface{} `json:"left_power_pct_avg"`
				RestActivity             interface{} `json:"RestActivity"`
			} `json:"Laps"`
		} `json:"Activities"`
		WarmUp           interface{} `json:"warm_up"`
		CoolDown         interface{} `json:"cool_down"`
		PlanInstanceInfo interface{} `json:"plan_instance_info"`
		IntegrationInfo  interface{} `json:"integration_info"`
	} `json:"data"`
	HideAfter        interface{} `json:"hide_after"`
	UserCurrentDate  interface{} `json:"user_current_date"`
	Success          bool        `json:"success"`
	ErrorNumber      interface{} `json:"error_number"`
	ErrorDescription interface{} `json:"error_description"`
	CallID           string      `json:"call_id"`
}