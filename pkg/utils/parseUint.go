package utils

import (
	"errors"
	"strconv"
)

var errIncorrectUintValue = errors.New("incorrect uint value")

func ParseUint(value string) (int32, error) {
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, errIncorrectUintValue
	}
	return int32(val), nil
}
