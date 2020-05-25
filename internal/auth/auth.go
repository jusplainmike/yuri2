package auth

import (
	"errors"
	"time"
)

var ErrInvalidSessionKey = errors.New("invalid session key")

type Options struct {
	ExpireTime time.Duration
}

type Auth interface {
	GenerateSessionKey(props interface{}) (string, error)
	ValidateSessionKey(key string) (interface{}, error)
	GetExpireTime() time.Duration
}
