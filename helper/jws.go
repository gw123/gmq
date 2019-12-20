package helper

import (
	"crypto/rsa"
	"errors"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/dgrijalva/jwt-go"
)

func ParseJwsTokenSh1(data string, secret string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(data, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*jwt.StandardClaims); ok {
		return c, nil
	}
	return nil, errors.New("转换标准jwt失败")
}

func MakeJwtTokenSh1(claims *jwt.StandardClaims, secret string) (string, error) {
	if claims == nil {
		return "", errors.New("MakeJwtToken claims  is nil")
	}

	mySigningKey := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ParseJwsTokenRs512(token []byte, pubkey *rsa.PublicKey) (map[string]interface{}, error) {
	jwstoken, err := jws.Parse(token)
	if err != nil {
		return nil, err
	}
	err = jwstoken.Verify(pubkey, crypto.SigningMethodRS512)
	if err != nil {
		return nil, err
	}
	maps, ok := jwstoken.Payload().(map[string]interface{})
	if !ok {
		return nil, errors.New("转换失败")
	}
	return maps, nil
}
