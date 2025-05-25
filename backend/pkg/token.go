package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TokenIssuer = "ZEIBA_APP"

type JWTMaker struct {
	secretKey string
}

// Payload defines the JWT payload structure
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uint32    `json:"user_id"`
	UserEmail string    `json:"user_email"`
	IsAdmin   bool      `json:"is_admin"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey}
}

func (maker *JWTMaker) CreateToken(
	userID uint32,
	email string,
	isAdmin bool,
	duration time.Duration,
) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	id, err := uuid.NewUUID()
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "failed to create uuid: %v", err)
	}

	payload := &Payload{
		ID:        id,
		UserID:    userID,
		UserEmail: email,
		IsAdmin:   isAdmin,
		IssuedAt:  now,
		ExpiresAt: exp,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "error signing token: %v", err)
	}
	return signedToken, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Payload{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, Errorf(INTERNAL_ERROR, "invalid signing method")
			}
			return []byte(maker.secretKey), nil
		},
	)
	if err != nil {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token: %v", err)
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, Errorf(INTERNAL_ERROR, "invalid token payload")
	}

	if payload.RegisteredClaims.Issuer != TokenIssuer {
		return nil, Errorf(INTERNAL_ERROR, "invalid issuer")
	}

	if payload.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, Errorf(INTERNAL_ERROR, "token is expired")
	}

	return payload, nil
}
