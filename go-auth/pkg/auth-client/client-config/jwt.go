package client_config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"username"` //账户
}

type Jwt struct {
	secret              string
	expireInHours       int64
	maxRefreshInMinutes int64
}

type JwtOption struct {
	Realm               string `mapstructure:"realm" json:"realm"`
	Secret              string `mapstructure:"secret" json:"secret"`
	ExpireTimeInHours   int64  `mapstructure:"expire-time-in-hours" json:"expireTimeInHours"`
	MaxRefreshInMinutes int64  `mapstructure:"max-refresh-in-minutes" json:"maxRefreshInMinutes"`
}

func NewJwt(secret string, expireInHours int64, maxRefreshMinutes int64) *Jwt {
	jwt := &Jwt{
		secret:              secret,
		expireInHours:       expireInHours, //小时为单位
		maxRefreshInMinutes: maxRefreshMinutes,
	}
	return jwt
}

func getKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
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

func GenPublic(privateKeyStr string) string {
	key, err := getKey(privateKeyStr)
	if err != nil {
		return ""
	}
	publicK := elliptic.Marshal(key.PublicKey.Curve, key.PublicKey.X, key.PublicKey.Y)
	publicStr := base64.StdEncoding.EncodeToString(publicK)
	//fmt.Printf("Secret: %s\n", publicStr)
	return publicStr
}

func (j *Jwt) CreateJwt(uid uint, username string) (string, error) {
	key, err := getKey(j.secret)
	//过期时间
	expiresAt := time.Now().Add(time.Hour * time.Duration(j.expireInHours)).Unix()
	//设置载荷
	claims := MyClaims{}
	claims.Id = fmt.Sprintf("%d", uid)
	claims.Username = username
	claims.ExpiresAt = expiresAt
	claims.Issuer = "xlnian@163.com" // 非必须，也可以填充用户名
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	//令牌签名
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParserToken 解析token
func (j *Jwt) ParserToken(tokenString string) (*MyClaims, error) {
	key, err := getKey(j.secret)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return &key.PublicKey, nil
	})

	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
			return nil, errors.New("token 已过期，请重新登录")
		}
		return claims, nil
	}
	return nil, errors.New("这不是一个 Token，请重新登录")
}

// ParserToken 解析token
func (j *Jwt) ValidAndRefreshToken(tokenString string) (*MyClaims, string, error) {
	claims, err := j.ParserToken(tokenString)
	if err != nil {
		return claims, "", err
	}

	if (claims.ExpiresAt - j.maxRefreshInMinutes) < time.Now().Unix() {
		//claims, tokenString, _ = j.CreateJwt(claims.Id, claims.Username, claims.CurrentTenantId, claims.UserType, claims.IsAdmin)
		//暂不刷新，考虑的问题太多了
		return claims, "", nil
	}

	return claims, "", nil
}
