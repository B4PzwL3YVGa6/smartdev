package util

import "github.com/anyandrea/smartdev/lib/env"

func GetUserAndPassword() (string, string) {
	return env.Get("WEATHERAPI_USERNAME", "anyandrea"), env.MustGet("WEATHERAPI_PASSWORD")
}

