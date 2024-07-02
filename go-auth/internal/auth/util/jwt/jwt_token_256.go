package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

func main11() {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
		return
	}

	//fmt.Printf("Secret: %s\n", base64.StdEncoding.EncodeToString(key))

	privateK, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("Secret: %s\n", base64.StdEncoding.EncodeToString(privateK))

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = 10
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenStr, err := t.SignedString(key)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("私钥签名的token: %s\n", tokenStr)

	// Validate token
	_, err1 := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return &key.PublicKey, nil
	})

	if err1 != nil {
		log.Fatalf("Token is invalid %v", err)
	} else {
		fmt.Println("Token is valid,publicKey %v:", &key.PublicKey)
	}
}
