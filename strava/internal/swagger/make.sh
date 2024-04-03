#!/bin/bash
# This will make the golang model code for the strava API using the swagger model provided by strava

swagger-codegen generate --input-spec https://developers.strava.com/swagger/swagger.json --lang go --output .