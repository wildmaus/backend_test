package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func ParseUintMass(r *http.Request, values ...string) ([]int32, error) {
	out := []int32{}
	for _, val := range values {
		v, err := strconv.ParseInt(mux.Vars(r)[val], 10, 32)
		if err != nil || v < 0 {
			return []int32{}, errIncorrectUintValue
		}
		out = append(out, int32(v))
	}
	return out, nil
}
