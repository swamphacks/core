package web

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	ErrNoPathParam = errors.New("path param is empty/missing")
)

func PathParamToUUID(r *http.Request, key string) (uuid.UUID, error) {
	s := chi.URLParam(r, key)
	if s == "" {
		return uuid.Nil, ErrNoPathParam
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
