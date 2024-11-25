package lib

import (
	"fmt"
	"time"

	"github.com/ramadhia/dataon-test/internal/config"
	"github.com/ramadhia/dataon-test/internal/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"user_id"`
	*jwt.RegisteredClaims
}

func GenerateTokens(claim model.Claim) (accessToken string, refreshToken string, err error) {
	cfg := config.Instance()
	accessTokenExpiry := time.Now().Add(1 * time.Hour)       // Access token berlaku 15 menit
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour) // Refresh token berlaku 7 hari

	accessClaims := &Claims{
		UserID: claim.ID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    cfg.App.Name,
			Subject:   claim.ID,
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
		},
	}

	refreshClaims := &Claims{
		UserID: claim.ID,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    cfg.App.Name,
			Subject:   claim.ID,
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiry),
		},
	}

	// Generate Access Token
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJwt.SignedString([]byte(cfg.App.JwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJwt.SignedString([]byte(cfg.App.JwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	cfg := config.Instance()

	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.App.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("refresh token tidak valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("gagal mendapatkan klaim dari refresh token")
	}

	return GenerateNewAccessToken(claims.UserID)
}

func GenerateNewAccessToken(userID string) (string, error) {
	cfg := config.Instance()

	accessTokenExpiry := time.Now().Add(15 * time.Minute)

	accessClaims := &Claims{
		UserID: userID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
		},
	}

	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return accessJwt.SignedString([]byte(cfg.App.JwtSecret))
}
