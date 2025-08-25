package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// return lowwer char
func Md5encoder(code string) string {
	m := md5.New()
	io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

func Md5StrToUpper(code string) string {
	return strings.ToUpper(Md5encoder(code))
}

func SaltPassWord(ps string, salt string) string {
	saltPw := fmt.Sprintf("%s$%s", Md5encoder(ps), salt)
	return saltPw
}

func CheckPassWord(rpw, salt, pw string) bool {
	return pw == SaltPassWord(rpw, salt)
}
