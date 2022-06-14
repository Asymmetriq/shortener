package models

import (
	"errors"
	"net/http"
)

var (
	ErrAlreadyExists = errors.New("such url is already shortened")
	ErrDeleted       = errors.New("short url deleted")
)

func ParsePostError(err error) int {
	switch {
	case err == nil:
		return http.StatusCreated
	case errors.Is(err, ErrAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}

}

func ParseGetError(err error) int {
	switch {
	case err == nil:
		return http.StatusTemporaryRedirect
	case errors.Is(err, ErrDeleted):
		return http.StatusGone
	default:
		return http.StatusBadRequest
	}
}
