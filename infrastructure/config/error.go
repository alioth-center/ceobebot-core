package config

import "fmt"

type NoMatchedKeyError struct {
	key string
}

func (e NoMatchedKeyError) Error() string {
	return fmt.Sprintf("no matched key [%s] in config content", e.key)
}

func NewNoMatchedKeyError(key string) error {
	return NoMatchedKeyError{key: key}
}
