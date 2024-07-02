package auth_client

import (
	"encoding/base64"
	"fmt"
	auth_jwt "gitee.com/go-nianxi/go-auth/pkg/auth-client/client-config"
	"strings"
	"testing"
)

func Test_connect(t *testing.T) {
	path := "D:\\favorite-work\\gostudy\\fabric-chaincode\\go-nianxi\\nianxi-auth\\configs\\auth.yml"
	authClient, _ := InitNxAuthClientByConfigFile(path)

	authClient.apiImpl.AuthApi.TryingConnect()

}

func Test_JWT(t *testing.T) {
	privateK := "MHcCAQEEII7ENwExt67UjqohHxUL5aNIsaNnUnjA7KJMqCueJLlzoAoGCCqGSM49AwEHoUQDQgAErE3pZstMaecYt0lnIMbDTDBn4NQYT05gkDIelhOgW7r1M3nx68VZ6yzTt3whDEjxYbyrNRQLlz8Yga0Qzff8SQ=="
	jwt := auth_jwt.NewJwt(privateK, 24, 15)

	//func (j *Jwt) CreateJwt(uid string, username string, currentTenantId uint, userType int8, isAdmin bool) (*MyClaims, string, error) {
	tokenStr, err := jwt.CreateJwt(1, "a")
	if err != nil {
		panic(err)
	}
	fmt.Sprintln(tokenStr)

	_, token, err1 := jwt.ValidAndRefreshToken(tokenStr)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(token)
	if err1 != nil {
		fmt.Println(err1.Error())
	} else {
		fmt.Println("token is valid!")
	}

}

func Test_DecodeString_JWT(t *testing.T) {
	privateK := "MHcCAQEEII7ENwExt67UjqohHxUL5aNIsaNnUnjA7KJMqCueJLlzoAoGCCqGSM49AwEHoUQDQgAErE3pZstMaecYt0lnIMbDTDBn4NQYT05gkDIelhOgW7r1M3nx68VZ6yzTt3whDEjxYbyrNRQLlz8Yga0Qzff8SQ=="
	jwt := auth_jwt.NewJwt(privateK, 24, 15)

	queryToken := "ZXlKaGJHY2lPaUpGVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk9EQTBNekF5TURJc0ltcDBhU0k2SWpFaUxDSnBjM01pT2lKdWFXRnVlR2tpTENKMWMyVnlibUZ0WlNJNkltRmtiV2x1SWl3aWFYTkJaRzFwYmlJNlptRnNjMlVzSW5WelpYSlVlWEJsSWpvd0xDSmpkWEp5Wlc1MFZHVnVZVzUwU1dRaU9qQjkuS3ByaUVrQU9kT0JuTUQzT3FPZnJOckVlbHdGLUNITllQcjZjcktsRUFtZkdZbHFaNEFydjNmak9ucFpFRHdyTVhCX0kxYUd1dHNyMnZWeThuazdvelE="
	tokenByte, _ := base64.URLEncoding.DecodeString(queryToken)
	queryToken = string(tokenByte)

	_, token, err1 := jwt.ValidAndRefreshToken(queryToken)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(token)
	if err1 != nil {
		fmt.Println(err1.Error())
	} else {
		fmt.Println("token is valid!")
	}

}

func Test_URl(t *testing.T) {
	str := "http://localhost:8000/api-doc/index.html#ajdfkld?token=ZXlKaGJHY2lPaUpGVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk9EQTBOVFkwTlRBc0ltcDBhU0k2SWpFaUxDSnBjM01pT2lKdWFXRnVlR2tpTENKMWMyVnlibUZ0WlNJNkltRmtiV2x1SWl3aWFYTkJaRzFwYmlJNlptRnNjMlVzSW5WelpYSlVlWEJsSWpvd0xDSmpkWEp5Wlc1MFZHVnVZVzUwU1dRaU9qQjkuZWtSSlduWmVJMzF1cFBCTFM5NGp2aDFwZE0tM2wwMTVxYWZRSUpjNl91c0ZRb2FGMHNMYVNVZGFJZTlPTW50a1FlMFA2WWxXaHVUUTZObllqREVTcnc=&sign=Qkt4TjZXYkxUR25uR0xkSlp5REd3MHd3WitEVUdFOU9ZSkF5SHBZVG9GdTY5VE41OGV2Rldlc3MwN2Q4SVF4SThXRzhxelVVQzVjL0dJR3RFTTMzL0VrPQ==&flag=bG9naW4="

	url := []string{""}
	url = append(url, substr(str, 0, strings.Index(str, "?")))
	fmt.Println(url)
	str = substr(str, strings.Index(str, "?")+1, len(str))
	list := strings.Split(str, "&")
	for _, v := range list {
		if !strings.HasPrefix(v, "sign=") && !strings.HasPrefix(v, "token=") && !strings.HasPrefix(v, "flag") {
			if len(url) == 1 {
				url = append(url, "?"+v)
			} else {
				url = append(url, "&"+v)
			}
		}
	}
	fmt.Println(url)

}
