package middleware

import (
	"errors"
	"fmt"
	"github.com/ramadhia/dataon-test/internal/model"
	"github.com/shortlyst-ai/go-helper"
	"net/http"
	"strings"

	"github.com/ramadhia/dataon-test/internal/handler/http/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	jwtData     = "JWT_DATA"
	tokenBearer = "TOKEN_BEARER"
)

type JWTData struct {
	jwt.RegisteredClaims
	model.Claim
}

func NewHmacJwtMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := getBearerAuth(c.Request)
		if bearer == nil {
			err := errors.New("missing bearer token")
			response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		jwtObj, err := decodeHmacJwtData(secretKey, *bearer)
		if err != nil {
			response.SendErrorResponse(c, response.ErrUnauthorized, "Unauthenticated")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if jwtObj.Claim.ID == "" {
			err := errors.New("invalid token")
			response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		fmt.Println(helper.MustJsonString(jwtObj))
		c.Set(jwtData, jwtObj.Claim)
		c.Set(tokenBearer, *bearer)

		c.Next()
	}
}

func GetClaim(c *gin.Context) (model.Claim, error) {
	anyObj, exist := c.Get(jwtData)
	if !exist {
		return model.Claim{}, errors.New("user not found")
	}

	claimed, validType := anyObj.(model.Claim)
	if !validType {
		return model.Claim{}, errors.New("invalid user type")
	}

	tokenObj, exist := c.Get(tokenBearer)
	if !exist {
		return model.Claim{}, errors.New("missing token")
	}

	token, validType := tokenObj.(string)
	if !validType {
		return model.Claim{}, errors.New("invalid token type")
	}

	claimed.Token = token
	return claimed, nil
}

func getBearerAuth(r *http.Request) *string {
	authHeader := r.Header.Get("Authorization")
	authForm := r.Form.Get("code")
	if authHeader == "" && authForm == "" {
		return nil
	}
	token := authForm
	if authHeader != "" {
		s := strings.SplitN(authHeader, " ", 2)
		if (len(s) != 2 || strings.ToLower(s[0]) != "bearer") && token == "" {
			return nil
		}
		// Use authorization header token only if token type is bearer else query string access token would be returned
		if len(s) > 0 && strings.ToLower(s[0]) == "bearer" {
			token = s[1]
		}
	}
	return &token
}

func decodeHmacJwtData(hmacSecret []byte, tokenStr string) (*JWTData, error) {
	var claim JWTData

	secretFn := func(token *jwt.Token) (interface{}, error) {
		if _, validSignMethod := token.Method.(*jwt.SigningMethodHMAC); !validSignMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &claim, secretFn)
	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*JWTData); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("invalid token")
}
