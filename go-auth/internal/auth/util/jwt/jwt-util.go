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

// 1.生成私钥，并用base64机密存储
func GenPrivateKey() {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
		return
	}

	privateK, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("Secret: %s\n", base64.StdEncoding.EncodeToString(privateK))
}

func main() {

	//keys := GenPrivateKey()

	privateK := "MHcCAQEEII7ENwExt67UjqohHxUL5aNIsaNnUnjA7KJMqCueJLlzoAoGCCqGSM49AwEHoUQDQgAErE3pZstMaecYt0lnIMbDTDBn4NQYT05gkDIelhOgW7r1M3nx68VZ6yzTt3whDEjxYbyrNRQLlz8Yga0Qzff8SQ=="

	//crateToken(privateK)

	tokenStr := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA0Mjc1ODEsImp0aSI6IjEiLCJpc3MiOiJuaWFueGkiLCJ1c2VybmFtZSI6ImFkbWluIiwiaXNBZG1pbiI6ZmFsc2UsInVzZXJUeXBlIjowLCJjdXJyZW50VGVuYW50SWQiOjB9.yHz1gToOIhoF1eO24LHY7zUvxdLJinxn0P-RYlC6wpjtMmmA2p9w0oZwGaEIRW&token=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA0Mjg1NjYsImp0aSI6IjEiLCJpc3MiOiJuaWFueGkiLCJ1c2VybmFtZSI6ImFkbWluIiwiaXNBZG1pbiI6ZmFsc2UsInVzZXJUeXBlIjowLCJjdXJyZW50VGVuYW50SWQiOjB9.j1xviT-Z09pxL8ySswVS0WzJRTCgTjzqTuGvoCjW_CprFlAYioYx3jXO-tAjV3n6q7cN6GTF0F3QOtxfEeoFWQ"

	validToken(tokenStr, privateK)
}

// 2.将私钥转为key
func GetPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParseECPrivateKey(bytes)

	if err != nil {
		return nil, err
	}
	return key, nil
}

// 3.将私钥转为key创建token
func crateToken(privateKey string) {
	key, err := GetPrivateKey(privateKey)
	if err != nil {
		log.Fatalln(err)
		return
	}

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
	fmt.Printf("token: %s\n", tokenStr)

}

// 4.用私钥对应的公钥验签
func validToken(tokenStr string, privateKey string) {
	key, err := GetPrivateKey(privateKey)
	if err != nil {
		log.Fatalln(err)
		return
	}

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
