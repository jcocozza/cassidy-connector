package main

import (
	"fmt"

	"github.com/jcocozza/cassidy-connector/strava/auth"
)

func main() {
    auth.InitialAuthorization()

    token, err := auth.GetAccessTokenFromAuthorizationCode("*** THE TOKEN FROM BROWSER HERE ***")

    if err != nil {
        panic(err)
    }

    fmt.Println(token)
}