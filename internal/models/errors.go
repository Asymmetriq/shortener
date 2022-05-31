package models

import (
	"errors"
	"net/http"
)

var ErrAlreadyExists = errors.New("such url is already shortened")

func ParseStorageError(err error) int {
	switch {
	case err == nil:
		return http.StatusCreated
	case errors.Is(err, ErrAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}

}
