package usecase

import (
	"log"
	"time"

	"user-service/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TokenUseCase struct {
	JwtConfig *model.JWTConfig
	Log       *logrus.Logger
}

func NewTokenUseCase(jwtConfig *model.JWTConfig, log *logrus.Logger) *TokenUseCase {
	return &TokenUseCase{
		JwtConfig: jwtConfig,
		Log:       log,
	}
}

func (c *TokenUseCase) GenerateToken(id int, role string) (string, time.Time, error) {

	tokenID := uuid.New().String()
	claims := &model.UserClaims{
		ID:      id,
		Role:    role,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.JwtConfig.AccessExpiry)),
			Issuer:    c.JwtConfig.Issuer,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(c.JwtConfig.SecretKey))
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return "", time.Time{}, err
	}

	expiresAccessToken := claims.ExpiresAt.Time

	return accessTokenString, expiresAccessToken, nil
}

func (c *TokenUseCase) ValidateToken(tokenString string) (*model.Auth, error) {
	var auth *model.Auth

	log.Println("tokenString: ", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(c.JwtConfig.SecretKey), nil
	})
	if err != nil {
		c.Log.Warnf("Failed to validate token : %+v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok || !token.Valid {
		c.Log.Warnf("Failed to validate token : %+v", jwt.ErrInvalidKey)
		return nil, jwt.ErrInvalidKey
	}

	auth = &model.Auth{
		ID:   claims.ID,
		Role: claims.Role,
	}

	return auth, nil
}
