# Cassidy Connector

Cassidy Connector is a part of the larger Cassidy project.

The goal of this code is to facilitate connections between a whole bunch of activity platforms.

This will start with a basic strava integration and (hopefully) grow overtime to include more.

The ultimate goal is to allow users to entirely disconnect their data from external servers if they so desire.

As a warning, everything should be considered unstable until v1.0.0 is released.

## Strava

The first step in the project is to get some basic data connections to allow users to import their data.
This will also give us a better feel of the structure and shape of data that we are dealing with.

At least in spirit, this portion will be modeled after [stravalib](https://github.com/stravalib/stravalib)

## Final Surge

Currently, Final Surge does not expose any api for public use, so this is a backengineering.
As such, it can break at any time. Moreover, its functionality is limited as I have not figured out all the endpoints.
