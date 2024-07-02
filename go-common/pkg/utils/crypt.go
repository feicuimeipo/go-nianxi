package utils

import (
	"fmt"
	"github.com/xlcetc/cryptogm/sm/sm3"
	"golang.org/x/crypto/bcrypt"
)

// 国密sm3
func getSm3(cid1 string) string {
	data := []byte(cid1)
	return fmt.Sprintf("%x", sm3.SumSM3(data))
}

// 密码加密 使用自适应hash算法, 不可逆
func GenPasswd(passwd string) string {
	hashPasswd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hashPasswd)
}

// 通过比较两个字符串hash判断是否出自同一个明文
// hashPasswd 需要对比的密文
// passwd 明文
func ComparePasswd(hashPasswd string, passwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd)); err != nil {
		return err
	}
	return nil
}
