package auth

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zekroTJA/yuri2/pkg/discordoauth"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256

	ErrInvalidClaimsType = errors.New("invalid claims type")
)

type JWTAuth struct {
	key []byte
	opt *Options
}

func NewJWTAuth(key []byte, opt *Options) (a *JWTAuth, err error) {
	if key == nil || len(key) == 0 {
		if key, err = generateAuthSecret(); err != nil {
			return
		}
	} else if len(key) < 32 {
		err = errors.New("key length must be at least 256 bit")
		return
	}

	a = &JWTAuth{
		key: key,
		opt: opt,
	}

	return
}

func (a *JWTAuth) GenerateSessionKey(props interface{}) (string, error) {
	user, ok := props.(*discordoauth.UserModel)
	if !ok {
		return "", errors.New("invalid property type")
	}

	now := time.Now()
	key, err := jwt.NewWithClaims(jwtSigningMethod, jwt.StandardClaims{
		Subject:   user.ID,
		ExpiresAt: now.Add(a.opt.ExpireTime).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Add(-1 * time.Minute).Unix(),
	}).SignedString(a.key)

	return key, err
}

func (a *JWTAuth) ValidateSessionKey(key string) (interface{}, error) {
	tokenObj, err := jwt.Parse(key, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err == jwt.ErrInvalidKey || !tokenObj.Valid {
		return nil, ErrInvalidSessionKey
	} else if err != nil {
		return nil, err
	}

	mClaims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaimsType
	}

	userId, _ := mClaims["sub"].(string)

	return userId, nil
}

func (a *JWTAuth) GetExpireTime() time.Duration {
	return a.opt.ExpireTime
}

func generateAuthSecret() (key []byte, err error) {
	key = make([]byte, 32)
	_, err = rand.Read(key)
	return
}
