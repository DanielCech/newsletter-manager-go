package common

import (
	"encoding/json"
)

func UnmarshalJSON[T any](instance T, jsonBody string) T {
	_ = json.Unmarshal([]byte(jsonBody), &instance)
	return instance
}

func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func Assert(condition bool, description string) {
	if !condition {
		panic(description)
	}
}
