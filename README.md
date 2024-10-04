# Cassidy Connector

Cassidy Connector is a part of the larger Cassidy project.

The goal of this code is to facilitate connections between a whole bunch of activity platforms.
For the foreseeable future, this project will only support reading from platforms, not writing to them.
While this is purely to reduce the total scope of the project, I have found that reading is much more common then writing.

This will start with a basic strava integration and (hopefully) grow overtime to include more.

As a warning, everything should be considered unstable until v1.0.0 is released.

## Design Philosophy
Every platform gets a package (e.g. `strava`, `finalSurge`).
### App
Each package contains an `app` folder.
The app folder represents an instance of an app created by the user.
For example, I need to go to strava and create an application with them.
I will be given credentials (e.g. client secret). These are used to instantiate the app.

In each app package is an `api` folder. This will contain a stuct with the name convention `<Platform>API` (e.g. `StravaAPI`).

The idea is that the App and `<Platform>API` structs are _long lived_.
You instantiate them when your app spins up and then use the same `<Platform>API` struct instance to make all of your calls.
(Or, if you have several, you can distribute the load across several)

### API struct
This api struct is the main point of access for users that want to programatically use the platform's api.
I like to think of this api struct as a convience wrapper. It hides the details of things like a swagger implementation, or even just raw http requests.
The api struct contains all the methods users will need for iteracting with the api.
It handles authentication, rate limits and will make calls to lower level methods that return things from the platform's api.

### The lower level
Technically speaking, the API struct should not call the platform api's directly.
There should be a lower level implementation folder in the main package folder that handles this.
The API struct will call these methods to make api requests.
These lower levels calls are exposed to the developer (in the `App` struct), however none of the convience of rate limiting, authentication, etc is handled.

When possible, this lower level implementation should be done using swagger and automatic code gen.

### cmd
Each platform package also has a `cmd` folder that provides an implementation of a CLI tool that can be used for easy testing of methods.
The CLI is not intended for any kind of heavy use. It is merely for ad-hoc work and testing.

## Strava

The first step in the project is to get some basic data connections to allow users to import their data.
This will also give us a better feel of the structure and shape of data that we are dealing with.

At least in spirit, this portion will be modeled after [stravalib](https://github.com/stravalib/stravalib).

## Final Surge

Currently, Final Surge does not expose any api for public use, so this is a backengineering.
As such, it can break at any time. Moreover, its functionality is limited as I have not figured out all the endpoints.
