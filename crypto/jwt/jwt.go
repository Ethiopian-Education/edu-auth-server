package jwt_jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/crypto/parser"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

type JWTClaims struct {
	// AllowedRoles []string `json:"allowed_roles,omitempty"`
	// LoginMethod  string           `json:"login_method"`
	// Roles        []string         `json:"roles"`
	jwt.RegisteredClaims
	Aud          string           `json:"aud"`
	Exp          *jwt.NumericDate `json:"exp"`
	Iat          *jwt.NumericDate `json:"iat"`
	Issuer       string           `json:"iss"`
	SignUpMethod string           `json:"signup_method"`
	Nonce        string           `json:"nonce"`
	Scope        []string         `json:"scope,omitempty"`
	Subject      string           `json:"sub"`
	TokenType    string           `json:"token_type"`
	First_name   string           `json:"first_name"`
	Middle_name  string           `json:"middle_name"`
	Last_name    string           `json:"last_name"`
	Gender       string           `json:"gender,omitempty"`
	Phone_number string           `json:"phone_number,omitempty"`
	Email        string           `json:"email"`
	Picture      string           `json:"picture,omitempty"`
	BirthDate    string           `json:"birthdate,omitempty"`
	Metadata     interface{}      `json:"metadata,omitempty"`
}

// Active Handler ( func )
func Sign(claims *JWTClaims) (string, error) {
	key, err := parser.ParseRsaPrivateKeyFromPemStr(config.PRIVATE_KEY)
	if err != nil {
		return "", nil
	}
	// exp_tt := time.Now().Add(24 * time.Hour).UTC().Unix()
	// var id jwt.ClaimStrings = jwt.ClaimStrings{"8cacd89f-9d0d-4035-a1a3-b1a338bef411"}
	now := time.Now().UTC()
	exp_tt := jwt.NewNumericDate(now.Add(24 * time.Hour))

	claims.Aud = "ethio-edu@gmail.com"
	claims.Exp = exp_tt
	claims.Iat = jwt.NewNumericDate(now)
	claims.Issuer = claims.Subject

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		logrus.Error("Signing error : ", err)
		return "", nil
	}
	return token, nil

}

// VALIDATE AND DECODE THE TOKEN HERE
func Validate(token string) (*JWTClaims, error) {
	var err error
	key, err := parser.ParseRsaPublicKeyFromPemStr(config.PUBLIC_KEY)
	if err != nil {
		logrus.Error("Signing error : ", err)
		return nil, nil
	}

	var Claims jwt.Claims = &JWTClaims{}

	tok, err := jwt.ParseWithClaims(token, Claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return &JWTClaims{}, err
	}

	claims, ok := tok.Claims.(*JWTClaims)
	if !ok || !tok.Valid {
		return &JWTClaims{}, errors.New("invalid-token")
	}
	// if claims.Exp < jwt.NewNumericDate(time.Now()) {

	// }

	return claims, nil

}

// GET TOKEN FROM HEADER
func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func CSign(metadata interface{}) (string, error) {
	claims := make(jwt.MapClaims)

	key, err := parser.ParseRsaPrivateKeyFromPemStr(config.PRIVATE_KEY)
	if err != nil {
		return "", nil
	}

	exp_tt := time.Now().Add(24 * time.Hour).UTC().Unix()
	// diff := time.Now().Sub(time.Unix(exp_tt, 0))

	claims["sub"] = "8cacd89f-9d0d-4035-a1a3-b1a338bef411"
	claims["metadata"] = metadata
	claims["exp"] = exp_tt
	claims["iat"] = time.Now().UTC().Unix()
	claims["aud"] = "ethio-edu@gmail.com"

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		logrus.Error("Signing error : ", err)
		return "", nil
	}
	return token, nil
}
