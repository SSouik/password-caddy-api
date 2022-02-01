package config

import (
	"os"
	"strconv"
)

type ConfigValue struct {
	Value string
}

/********** PARSERS **********/

/**
Parse a string into a base 10 integer
*/
func ParseInt(value string) int64 {
	output, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return 0
	}

	return output
}

/**
Parse a string into a boolean value (true or false)
*/
func ParseBool(value string) bool {
	output, err := strconv.ParseBool(value)

	if err != nil {
		return false
	}

	return output
}

/********** CONVERTERS **********/

/**
Convert a ConfigValue into a int64
*/
func (value ConfigValue) ToInt64() int64 {
	return ParseInt(value.Value)
}

/**
Convert a ConfigValue into a string
*/
func (value ConfigValue) ToString() string {
	return value.Value
}

/**
Get an environment variable and return a ConfigValue.
If the environment variable does not exist, then return the default
value.
*/
func Get(key, defaultValue string) ConfigValue {
	var value = os.Getenv(key)

	if value == "" {
		return ConfigValue{
			Value: defaultValue,
		}
	}

	return ConfigValue{
		Value: value,
	}
}
