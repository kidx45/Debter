package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (maker *JWTMaker) CreateToken(userID int64, username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        payload.ID,
		"userId":    payload.UserID,
		"username":  payload.Username,
		"issuedAt":  payload.IssuedAt.Unix(),
		"expiredAt": payload.ExpiredAt.Unix(),
	})
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(fmt.Sprintf("%v", claims["id"]))
	if err != nil {
		return nil, ErrInvalidToken
	}
	userIDFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	issuedAtFloat, ok := claims["issuedAt"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	expiredAtFloat, ok := claims["expiredAt"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		ID:        id,
		UserID:    int64(userIDFloat),
		Username:  username,
		IssuedAt:  time.Unix(int64(issuedAtFloat), 0),
		ExpiredAt: time.Unix(int64(expiredAtFloat), 0),
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
