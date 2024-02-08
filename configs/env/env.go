package env

import (
	"os"

	"github.com/go-anon/simple/types/ints"
	"github.com/go-anon/simple/types/strs"
)

// Get
func Get(env_name, defaults string) string {
	env, ok := os.LookupEnv(env_name)
	if ok {
		return env
	}
	return defaults
}

// GetSlice
func GetSlice(env_name string, defaults ...string) []string {
	string_values := Get(env_name, "")
	if len(string_values) > 0 {
		return strs.SplitString(string_values, ",")
	}
	return defaults
}

// GetInt
func GetInt(env_name string, defaults int) int {
	env, ok := os.LookupEnv(env_name)
	if ok {
		return ints.Get(env)
	}
	return defaults
}

// GetSliceInt
func GetSliceInt(env_name string, defaults ...int) (res []int) {
	values := GetSlice(env_name)
	if len(values) > 0 {
		for _, v := range values {
			res = append(res, ints.Get(v))
		}
		return
	}
	return defaults
}
