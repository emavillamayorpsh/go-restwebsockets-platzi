package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId string `json:"userId`

	// with this line of code now "AppClaims" have all the properties defined inside of "jwt.StandardClaims"  (Audience, Id , ExpiresAt, etc)
	jwt.StandardClaims
}