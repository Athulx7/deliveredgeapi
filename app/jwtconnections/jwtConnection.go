package jwtconnections

import (
	"errors"
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/revel/revel"
)

var jwtKey []byte

func InitJWTSecret() {
	secret := revel.Config.StringDefault("jwt.secret", "")
	if secret == "" {
		revel.AppLog.Warn("⚠️  Missing jwt.secret in app.conf — using default key (NOT secure!)")
		fmt.Println("⚠️  Missing jwt.secret in app.conf — using default key (NOT secure!)")
		secret = "DefaultSecretDeliverEdgeKey"
	}
	jwtKey = []byte(secret)
}

type JWTClaims struct {
	UserID      int    `json:"user_id"`
	CompanyID   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
	UserCode    string `json:"user_code"`
	DBName      string `json:"db_name"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, companyID int, companyName, dbName string, UserCode string) (string, error) {
	fmt.Println("Generating JWT for UserID:", userID, "CompanyID:", companyID)
	if len(jwtKey) == 0 {
		InitJWTSecret()
	}

	claims := JWTClaims{
		UserID:      userID,
		CompanyID:   companyID,
		CompanyName: companyName,
		DBName:      dbName,
		UserCode:    UserCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "DeliverEdgeAPI",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		revel.AppLog.Errorf("❌ JWT signing failed: %v", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	if len(jwtKey) == 0 {
		InitJWTSecret()
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
